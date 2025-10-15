class PixelBattle {
  constructor(canvasId, text, font) {
    this.canvas = document.getElementById(canvasId);
    if (!this.canvas) {
      console.error(`Canvas with id "${canvasId}" not found.`);
      return;
    }
    this.ctx = this.canvas.getContext("2d");
    this.text = text;
    this.pixelSize = 8;
    this.font = font;
    this.color = "red";
    this.colorPicker = null;

    this.init();
  }

  init() {
    this.setupCanvas();
    this.createTextGrid();
    this.drawGrid();
    this.loadPixels();
    this.addEventListeners();
  }

  loadPixels() {
    fetch("/pixels")
      .then((response) => {
        if (!response.ok) {
          throw new Error("Failed to load pixels");
        }
        return response.json();
      })
      .then((pixels) => {
        pixels.forEach((pixel) => {
          if (this.textPixels[pixel.y] && this.textPixels[pixel.y][pixel.x]) {
            this.ctx.fillStyle = pixel.color;
            this.ctx.fillRect(pixel.x * this.pixelSize, pixel.y * this.pixelSize, this.pixelSize, this.pixelSize);
            this.ctx.strokeRect(pixel.x * this.pixelSize, pixel.y * this.pixelSize, this.pixelSize, this.pixelSize);
          }
        });
      })
      .catch((error) => {
        console.error("Error loading pixels:", error);
      });
  }

  setupCanvas() {
    const textCanvas = document.createElement("canvas");
    const textCtx = textCanvas.getContext("2d");
    textCtx.font = this.font;
    const textMetrics = textCtx.measureText(this.text);

    this.canvas.width = Math.ceil(textMetrics.width) + 2 * this.pixelSize;
    this.canvas.height = 300;
    this.canvas.style.width = `${this.canvas.width}px`;
    this.canvas.style.height = `${this.canvas.height}px`;

    textCanvas.width = this.canvas.width;
    textCanvas.height = this.canvas.height;
    textCtx.font = this.font;
    textCtx.fillStyle = "#000";
    textCtx.textAlign = "center";
    textCtx.textBaseline = "middle";
    textCtx.fillText(this.text, textCanvas.width / 2, textCanvas.height / 2);

    this.textCanvas = textCanvas;
    this.textCtx = textCtx;
  }

  createTextGrid() {
    const gridWidth = Math.floor(this.canvas.width / this.pixelSize);
    const gridHeight = Math.floor(this.canvas.height / this.pixelSize);
    const imageData = this.textCtx.getImageData(0, 0, this.textCanvas.width, this.textCanvas.height);
    const pixelData = imageData.data;
    this.textPixels = new Array(gridHeight).fill(null).map(() => new Array(gridWidth).fill(false));

    for (let y = 0; y < gridHeight; y++) {
      for (let x = 0; x < gridWidth; x++) {
        const pixelIndex = (y * this.pixelSize * this.textCanvas.width + x * this.pixelSize) * 4;
        if (pixelData[pixelIndex + 3] > 0) {
          this.textPixels[y][x] = true;
        }
      }
    }
  }

  drawGrid() {
    this.ctx.strokeStyle = "#ccc";
    const gridWidth = this.textPixels[0].length;
    const gridHeight = this.textPixels.length;

    for (let y = 0; y < gridHeight; y++) {
      for (let x = 0; x < gridWidth; x++) {
        if (this.textPixels[y][x]) {
          this.ctx.strokeRect(x * this.pixelSize, y * this.pixelSize, this.pixelSize, this.pixelSize);
        }
      }
    }
  }

  addEventListeners() {
    this.canvas.addEventListener("click", (e) => {
      const rect = this.canvas.getBoundingClientRect();
      const x = e.clientX - rect.left;
      const y = e.clientY - rect.top;

      const gridWidth = this.textPixels[0].length;
      const pixelX = Math.floor(x / (rect.width / gridWidth));
      const pixelY = Math.floor(y / (rect.height / this.textPixels.length));

      if (!this.textPixels[pixelY] || !this.textPixels[pixelY][pixelX]) {
        return;
      }

      const color = this.color;

      fetch("/pixels:paint", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          x: pixelX,
          y: pixelY,
          color: color,
        }),
      })
        .then((response) => {
          if (response.ok) {
            this.ctx.fillStyle = color;
            this.ctx.fillRect(pixelX * this.pixelSize, pixelY * this.pixelSize, this.pixelSize, this.pixelSize);
            this.ctx.strokeRect(pixelX * this.pixelSize, pixelY * this.pixelSize, this.pixelSize, this.pixelSize);
          } else {
            console.error("Failed to paint pixel");
          }
        })
        .catch((error) => {
          console.error("Error:", error);
        });
    });

    this.canvas.addEventListener("contextmenu", (e) => {
      e.preventDefault();
      this.showColorPicker(e.clientX, e.clientY);
    });
  }

  showColorPicker(x, y) {
    this.hideColorPicker(); // Hide any existing color picker

    const colorPicker = document.createElement("div");
    colorPicker.className = "color-picker";
    document.body.appendChild(colorPicker);
    this.colorPicker = colorPicker;

    const colors = ["red", "green", "blue", "yellow", "purple", "orange", "black", "white"];
    const radius = 50; // Radius of the circle
    const angleStep = (2 * Math.PI) / colors.length;

    colors.forEach((color, index) => {
      const colorOption = document.createElement("div");
      colorOption.className = "color-option";
      colorOption.style.backgroundColor = color;

      const angle = index * angleStep;
      const optionX = radius * Math.cos(angle);
      const optionY = radius * Math.sin(angle);

      colorOption.style.transform = `translate(${optionX}px, ${optionY}px)`;

      colorOption.addEventListener("click", () => {
        this.selectColor(color);
      });

      colorPicker.appendChild(colorOption);
    });

    colorPicker.style.left = `${x}px`;
    colorPicker.style.top = `${y}px`;

    // Add a listener to close the color picker when clicking outside
    setTimeout(() => {
      document.addEventListener("click", this.handleOutsideClick.bind(this));
    }, 0);
  }

  handleOutsideClick(e) {
    if (this.colorPicker && !this.colorPicker.contains(e.target)) {
      this.hideColorPicker();
    }
  }

  hideColorPicker() {
    if (this.colorPicker) {
      this.colorPicker.remove();
      this.colorPicker = null;
      document.removeEventListener("click", this.handleOutsideClick.bind(this));
    }
  }

  selectColor(color) {
    this.color = color;
    this.hideColorPicker();
  }
}

class Fish {
  constructor(ctx, canvasWidth) {
    this.ctx = ctx;
    this.canvasWidth = canvasWidth;
    this.isOffscreen = true;
  }

  respawn(seed) {
    const { fishData, palette, scale } = generateRandomFish(seed);
    this.fishData = fishData;
    this.palette = palette;
    this.pixelSize = 4 * scale;
    this.x = -this.getWidth() - Math.random() * 500;
    this.y = Math.random() * infiniteCanvas.height;
    this.speedX = Math.random() * 0.5 + 0.1;
    this.isOffscreen = false;

    this.offscreenCanvas = document.createElement("canvas");
    this.offscreenCanvas.width = this.getWidth();
    this.offscreenCanvas.height = this.getHeight();
    this.offscreenCtx = this.offscreenCanvas.getContext("2d");
    this.offscreenCtx.imageSmoothingEnabled = false;
    this.drawToOffscreenCanvas();
  }

  drawToOffscreenCanvas() {
    for (let r = 0; r < this.fishData.length; r++) {
      for (let c = 0; c < this.fishData[r].length; c++) {
        const colorIndex = this.fishData[r][c];
        if (colorIndex) {
          this.offscreenCtx.fillStyle = this.palette[colorIndex];
          this.offscreenCtx.fillRect(c * this.pixelSize, r * this.pixelSize, this.pixelSize, this.pixelSize);
        }
      }
    }
  }

  update() {
    if (this.isOffscreen) return;
    this.x += this.speedX;
    if (this.x > this.canvasWidth) {
      this.isOffscreen = true;
    }
  }

  draw() {
    if (this.isOffscreen) return;
    this.ctx.drawImage(this.offscreenCanvas, this.x, this.y);
  }

  getWidth() {
    return this.fishData[0].length * this.pixelSize;
  }

  getHeight() {
    return this.fishData.length * this.pixelSize;
  }
}

class FishManager {
  constructor(ctx, canvasWidth) {
    this.ctx = ctx;
    this.canvasWidth = canvasWidth;
    this.fishes = [];
    this.seedBuffer = [];
    this.currentPage = 1;
    this.isLoading = false;
    this.TARGET_FISH_COUNT = 50;
    this.BUFFER_REFILL_THRESHOLD = 10;
    this.ADD_FISH_INTERVAL = 1000; // Add a new fish every second
  }

  async fillSeedBuffer() {
    if (this.isLoading || this.seedBuffer.length > this.BUFFER_REFILL_THRESHOLD) {
      return;
    }
    this.isLoading = true;
    try {
      const response = await fetch(`/fishes?page=${this.currentPage}`);
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const data = await response.json();
      this.seedBuffer.push(...data);
      if (data.length < 100) {
        this.currentPage = 1;
      } else {
        this.currentPage++;
      }
    } catch (error) {
      console.error("Error loading fish seeds:", error);
      setTimeout(() => {
        this.isLoading = false;
        this.fillSeedBuffer();
      }, 5000);
      return; // Return here to avoid setting isLoading to false immediately
    }
    this.isLoading = false;
  }

  addFish() {
    if (this.fishes.length >= this.TARGET_FISH_COUNT) {
      return;
    }
    const seed = this.seedBuffer.shift();
    if (seed) {
      const fish = new Fish(this.ctx, this.canvasWidth);
      fish.respawn(seed);
      this.fishes.push(fish);
    }
    this.fillSeedBuffer();
  }

  start() {
    this.fillSeedBuffer();
    setInterval(() => this.addFish(), this.ADD_FISH_INTERVAL);
    this.animate();
  }

  animate() {
    this.ctx.clearRect(0, 0, this.canvasWidth, infiniteCanvas.height);
    this.fishes.forEach((fish) => {
      fish.update();
      fish.draw();
      if (fish.isOffscreen) {
        const seed = this.seedBuffer.shift();
        if (seed) {
          fish.respawn(seed);
        }
        this.fillSeedBuffer();
      }
    });
    requestAnimationFrame(() => this.animate());
  }
}

const emailLink = document.getElementById("email-link-handler");
emailLink.addEventListener("click", (e) => {
  e.preventDefault();
  const user = "me";
  const domain = "tomashevich";
  window.location.href = `mailto:${user}@${domain}`;
});

new PixelBattle("pixel-canvas", "tomashevich", "bold 200px Lato");

const infiniteCanvas = document.getElementById("infinite-canvas");
const infiniteCtx = infiniteCanvas.getContext("2d");
infiniteCtx.imageSmoothingEnabled = false;

const bubbleCanvas = document.getElementById("bubble-canvas");
const bubbleCtx = bubbleCanvas.getContext("2d");

function resizeCanvas() {
  infiniteCanvas.width = window.innerWidth;
  infiniteCanvas.height = window.innerHeight;
  bubbleCanvas.width = window.innerWidth;
  bubbleCanvas.height = window.innerHeight;
}

window.addEventListener("resize", resizeCanvas);
resizeCanvas();

function generateRandomFish(seed) {
  class SeededRandom {
    constructor(seed) {
      this.seed = seed;
    }

    next() {
      this.seed = (this.seed * 9301 + 49297) % 233280;
      return this.seed / 233280;
    }
  }

  let seedValue = 0;
  for (let i = 0; i < seed.length; i++) {
    seedValue += seed.charCodeAt(i);
  }
  const random = new SeededRandom(seedValue);

  function generateFishBody(width, height) {
    const body = Array(height)
      .fill(0)
      .map(() => Array(width).fill(0));
    const a = width / 2;
    const b = height / 2;
    for (let y = 0; y < height; y++) {
      for (let x = 0; x < width; x++) {
        if ((x - a) ** 2 / a ** 2 + (y - b) ** 2 / b ** 2 < 1) {
          body[y][x] = 1;
        }
      }
    }
    return body;
  }

  function generateFishTail(width, height, type) {
    const tail = Array(height)
      .fill(0)
      .map(() => Array(width).fill(0));
    switch (type) {
      case 0:
        for (let y = 0; y < height; y++) {
          for (let x = 0; x < width; x++) {
            if (y >= x && y <= height - x) {
              tail[y][x] = 1;
            }
          }
        }
        break;
      case 1: // Forked tail
        for (let y = 0; y < height; y++) {
          for (let x = 0; x < width; x++) {
            if (y >= x && y <= height - x && y > height / 2 - 1 && y < height / 2 + 1) {
              tail[y][x] = 0;
            } else if (y >= x && y <= height - x) {
              tail[y][x] = 1;
            }
          }
        }
        break;
      case 2: // Long tail
        for (let y = 0; y < height; y++) {
          for (let x = 0; x < width; x++) {
            if (y > height / 2 - 2 && y < height / 2 + 2) {
              tail[y][x] = 1;
            }
          }
        }
        break;
      // Add 7 more tail types
      default:
        for (let y = 0; y < height; y++) {
          for (let x = 0; x < width; x++) {
            if (y >= x && y <= height - x) {
              tail[y][x] = 1;
            }
          }
        }
        break;
    }
    return tail;
  }

  function addBodyPatterns(body) {
    const patternType = Math.floor(random.next() * 4);
    for (let y = 0; y < body.length; y++) {
      for (let x = 0; x < body[y].length; x++) {
        if (body[y][x] === 1) {
          switch (patternType) {
            case 0: // Spots
              if (random.next() > 0.7) {
                body[y][x] = 2;
              }
              break;
            case 1: // Wavy lines
              if (Math.sin(x * 0.5 + y * 0.5) > 0.5) {
                body[y][x] = 2;
              }
              break;
            case 2: // Lines
              if (y % 3 === 0) {
                body[y][x] = 2;
              }
              break;
            case 3: // Crosses
              if (x % 4 === 0 || y % 4 === 0) {
                body[y][x] = 2;
              }
              break;
          }
        }
      }
    }
    return body;
  }

  function addAnglerfishLight(head) {
    const lightX = Math.floor(head[0].length / 2);
    for (let y = 0; y < head.length / 2; y++) {
      head[y][lightX] = 5; // Light color index
    }
    head[0][lightX] = 0;
    head[1][lightX - 1] = 5;
    head[1][lightX + 1] = 5;
    return head;
  }

  function generateProceduralPalette() {
    const baseHue = random.next() * 360;
    const palette = {
      1: `hsl(${baseHue}, 70%, 50%)`,
      2: `hsl(${(baseHue + 120) % 360}, 70%, 60%)`,
      3: "#000000",
      4: "#ffffff",
      5: "#ffff00", // Light color
    };
    return palette;
  }

  function combineParts(parts) {
    const tail = parts.tail;
    const body = parts.body;
    const head = parts.head;

    const fishWidth = tail[0].length + body[0].length + head[0].length - 4; // Overlap by 2 pixels on each side
    const fishHeight = body.length;
    const fishData = Array(fishHeight)
      .fill(0)
      .map(() => Array(fishWidth).fill(0));

    // Draw tail
    for (let y = 0; y < tail.length; y++) {
      for (let x = 0; x < tail[y].length; x++) {
        if (tail[y][x]) {
          fishData[y][x] = tail[y][x];
        }
      }
    }

    // Draw body
    let currentX = tail[0].length - 2;
    for (let y = 0; y < body.length; y++) {
      for (let x = 0; x < body[y].length; x++) {
        if (body[y][x]) {
          fishData[y][currentX + x] = body[y][x];
        }
      }
    }

    // Draw head
    currentX += body[0].length - 2;
    for (let y = 0; y < head.length; y++) {
      for (let x = 0; x < head[y].length; x++) {
        if (head[y][x]) {
          fishData[y][currentX + x] = head[y][x];
        }
      }
    }

    return fishData;
  }

  const hasTeeth = random.next() > 0.8;
  const hasAnglerfishLight = random.next() > 0.9;
  const hasPatterns = random.next() > 0.5;

  const headHeight = Math.floor(random.next() * 4) + 8;
  const bodyHeight = headHeight;
  const tailHeight = headHeight;

  const headWidth = Math.floor(random.next() * 2) + 6;
  const bodyWidth = Math.floor(random.next() * 8) + 10;
  const tailWidth = Math.floor(random.next() * 3) + 5;
  const tailType = Math.floor(random.next() * 10);

  let body = generateFishBody(bodyWidth, bodyHeight);
  if (hasPatterns) {
    body = addBodyPatterns(body);
  }

  let head = generateFishBody(headWidth, headHeight);
  if (hasAnglerfishLight) {
    head = addAnglerfishLight(head);
  }

  const parts = {
    tail: generateFishTail(tailWidth, tailHeight, tailType),
    body: body,
    head: head,
  };

  // Add eyes
  const eyeY = Math.floor(headHeight / 2);
  const eyeX = Math.floor(headWidth - 3);
  parts.head[eyeY][eyeX] = 3;

  // Add teeth
  if (hasTeeth) {
    for (let i = 0; i < headWidth; i++) {
      if (i > headWidth / 2 && random.next() > 0.5) {
        parts.head[headHeight - 2][i] = 4;
      }
    }
  }

  const fishData = combineParts(parts);
  const palette = generateProceduralPalette();
  const scale = random.next() * 0.3 + 0.3;

  return { fishData, palette, scale };
}

const fishManager = new FishManager(infiniteCtx, infiniteCanvas.width);
fishManager.start();

class Bubble {
  constructor() {
    this.x = Math.random() * bubbleCanvas.width;
    this.y = bubbleCanvas.height + Math.random() * 100;
    this.radius = Math.random() * 3 + 1;
    this.speedY = Math.random() * 0.5 + 0.2;
  }

  update() {
    this.y -= this.speedY;
  }

  draw() {
    bubbleCtx.beginPath();
    bubbleCtx.arc(this.x, this.y, this.radius, 0, Math.PI * 2);
    bubbleCtx.fillStyle = "rgba(193, 236, 250, 0.9)";
    bubbleCtx.fill();
  }
}

const bubbles = [];

function animateBubbles() {
  bubbleCtx.clearRect(0, 0, bubbleCanvas.width, bubbleCanvas.height);

  if (Math.random() < 0.1) {
    bubbles.push(new Bubble());
  }

  for (let i = bubbles.length - 1; i >= 0; i--) {
    bubbles[i].update();
    bubbles[i].draw();

    if (bubbles[i].y < -bubbles[i].radius) {
      bubbles.splice(i, 1);
    }
  }

  requestAnimationFrame(animateBubbles);
}

animateBubbles();
