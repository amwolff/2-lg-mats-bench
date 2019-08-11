#ifndef SERVER_MATMULT_IMPL_H
#define SERVER_MATMULT_IMPL_H

#include "amwolff/matmult/v1/matrix_product_api.grpc.pb.h"

using amwolff::matmult::v1::MatrixProductAPI;
using amwolff::matmult::v1::MultiplyRequest;
using amwolff::matmult::v1::MultiplyResponse;
using grpc::ServerContext;
using grpc::Status;

class MatrixProductAPIImpl final : public MatrixProductAPI::Service {
  Status Multiply(ServerContext* context, const MultiplyRequest* request,
                  MultiplyResponse* response) override;
};

#endif  // SERVER_MATMULT_IMPL_H
