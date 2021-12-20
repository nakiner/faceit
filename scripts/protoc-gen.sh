docker run --platform linux/amd64 --rm -v $(pwd):/defs namely/protoc-all \
    -d api/protobuf-spec \
    -i scripts \
    -i vendor \
    -o . \
    -l go