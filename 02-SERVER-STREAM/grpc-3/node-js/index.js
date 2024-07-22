const grpc = require('@grpc/grpc-js');
const protoLoader = require('@grpc/proto-loader');

const SERVER = 'localhost:50051';

// Load the protobuf files
const protoFiles = ["../proto/product_model.proto", "../proto/product_service.proto"];
const packageDefinition = protoLoader.loadSync(protoFiles, {
    keepCase: true,
    longs: String,
    enums: String,
    defaults: true,
    oneofs: true
});

const protoDescriptor = grpc.loadPackageDefinition(packageDefinition);
const productProto = protoDescriptor.product; // Assuming the package name in your .proto files is 'product'

const client = new productProto.ProductService(SERVER, grpc.credentials.createInsecure());

function getSingleProduct() {
    const request = {}; // Empty request for ProductRequest

    client.GetProduct(request, (error, response) => {
        if (error) {
            console.error("Error:", error);
            return;
        }

        console.log("Response from the server:");
        response.products.forEach((product) => {
            console.log(`id: ${product.id}, title: ${product.title}`);
        });
    });
}

// Start the client
getSingleProduct();
