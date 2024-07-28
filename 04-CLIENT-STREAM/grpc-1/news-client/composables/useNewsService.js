// import { ref } from '@nuxtjs/composition-api';
import grpc from '@grpc/grpc-js';
import protoLoader from '@grpc/proto-loader';

export const useNewsService = () => {
  const news = ref([]);
  const error = ref(null);

  const getNewsBulk = async () => {
    const SERVER = 'localhost:50051';

    // Load the protobuf files
    const protoFiles = ["../proto/news_models.proto", "../proto/news_service.proto"];
    const packageDefinition = protoLoader.loadSync(protoFiles, {
        keepCase: true,
        longs: String,
        enums: String,
        defaults: true,
        oneofs: true
    });

    const protoDescriptor = grpc.loadPackageDefinition(packageDefinition);
    const newsProto = protoDescriptor.news; // Assuming the package name in your .proto files is 'news'

    const client = new newsProto.Newservice(SERVER, grpc.credentials.createInsecure());

    const request = {}; // Empty request for NewsRequest

    client.GetNewsBulk(request, (err, response) => {
        if (err) {
            error.value = err;
            return;
        }

        news.value = response.news;
    });
  };

  return { news, error, getNewsBulk };
}
