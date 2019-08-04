#include <iostream>
#include <string>

#include <grpcpp/grpcpp.h>

#include "perf_srvc_impl.h"

using grpc::Server;
using grpc::ServerBuilder;

void StartServer(const std::string &server_address = "0.0.0.0:50051") {
  ServerBuilder builder;

  builder.AddListeningPort(server_address, grpc::InsecureServerCredentials());

  PerformerServiceImpl service;
  builder.RegisterService(&service);

  const std::unique_ptr<Server> server(builder.BuildAndStart());
  std::cout << "Server listening on " << server_address << std::endl;

  server->Wait();
}

int main(int argc, char **argv) {
  StartServer();
  return EXIT_SUCCESS;
}
