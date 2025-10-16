package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/zstd"
)

type CompressResponseWriter struct {
	http.ResponseWriter
	writer io.WriteCloser
}

func (w CompressResponseWriter) Write(data []byte) (int, error) {
	return w.writer.Write(data)
}

func (w CompressResponseWriter) WriteHeader(code int) {
	w.ResponseWriter.Header().Del("Content-Length")
	w.ResponseWriter.WriteHeader(code)
}

func Compress() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			acceptEncoding := r.Header.Get("Accept-Encoding")

			if strings.Contains(acceptEncoding, "br") {
				brotliWriter := brotli.NewWriterV2(w, 9)
				w.Header().Set("Content-Encoding", "br")
				defer brotliWriter.Close()
				next.ServeHTTP(CompressResponseWriter{ResponseWriter: w, writer: brotliWriter}, r)
				return
			}

			if strings.Contains(acceptEncoding, "zstd") {
				enc, err := zstd.NewWriter(w, zstd.WithEncoderLevel(zstd.SpeedBestCompression))
				if err == nil {
					w.Header().Set("Content-Encoding", "zstd")
					defer enc.Close()
					next.ServeHTTP(CompressResponseWriter{ResponseWriter: w, writer: enc}, r)
					return
				}
			}

			if strings.Contains(acceptEncoding, "gzip") {
				gz := gzip.NewWriter(w)
				w.Header().Set("Content-Encoding", "gzip")
				defer gz.Close()
				next.ServeHTTP(CompressResponseWriter{ResponseWriter: w, writer: gz}, r)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
