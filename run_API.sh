echo "Delete existing Likes container"

docker stop neo4j PerfectPick_Likes_ms
docker rm neo4j PerfectPick_Likes_ms

echo "Delete existing Likes image"

docker rmi neo4j PerfectPick_Likes_ms

echo "running docker..."

docker-compose up