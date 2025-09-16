# ec2-builder

go build -o ec2-builder

./ec2-builder --help
./ec2-builder create --type "t3.xlarge" --name "myproject" --image "ami-087f82ea9fe6a26e8"
./ec2-builder list