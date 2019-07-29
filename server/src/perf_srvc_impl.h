#ifndef OPSERVER_PERF_SRVC_IMPL_H
#define OPSERVER_PERF_SRVC_IMPL_H

#include "matmult.grpc.pb.h"

using grpc::ServerContext;
using grpc::Status;
using grpc::StatusCode;
using matmult::Performer;
using matmult::PerformerRequest;
using matmult::PerformerResponse;

class PerformerServiceImpl final : public Performer::Service {
  Status MultiplyMatrices(ServerContext* context,
                          const PerformerRequest* request,
                          PerformerResponse* response) override;
};

#endif  // OPSERVER_PERF_SRVC_IMPL_H
