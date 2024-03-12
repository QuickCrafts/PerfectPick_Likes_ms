
echo "docker build..."
docker build . -t go-containerized:latest

echo ""
echo ""
echo "docker run..."
docker run -p 3000:3000 go-containerized:latest