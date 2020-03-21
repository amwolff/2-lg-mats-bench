#include <iostream>
#include <string>

#include <grpcpp/ext/proto_server_reflection_plugin.h>
#include <grpcpp/grpcpp.h>
#include <Eigen/Core>

#include "matmult_impl.h"

using grpc::Server;
using grpc::ServerBuilder;

void StartServer(const std::string &server_address = "0.0.0.0:50051") {
  grpc::EnableDefaultHealthCheckService(true);
  grpc::reflection::InitProtoReflectionServerBuilderPlugin();

  ServerBuilder builder;

  builder.AddListeningPort(server_address, grpc::InsecureServerCredentials());

  MatrixProductAPIImpl service;
  builder.RegisterService(&service);
  builder.SetMaxReceiveMessageSize(INT_MAX);

  std::unique_ptr<Server> server(builder.BuildAndStart());
  std::cout << "Server listening on " << server_address << std::endl;

  server->Wait();
}

int main(int argc, char **argv) {
  const int n = Eigen::nbThreads();
  if (n > 1) {
    std::cout << "Eigen: parallelizing on " << n << " threads" << std::endl;
  }

  StartServer();

  return EXIT_SUCCESS;
}
