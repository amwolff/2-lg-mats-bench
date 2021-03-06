#include "matmult_impl.h"

#include <Eigen/Dense>

using Eigen::Index;
using Eigen::MatrixXd;
using grpc::StatusCode;

Status MatrixProductAPIImpl::Multiply(ServerContext* context,
                                      const MultiplyRequest* request,
                                      MultiplyResponse* response) {
  const auto& mplier = request->multiplier();
  const int mplier_columns_size = mplier.columns_size();
  const int mplier_rows_size = mplier.columns(0).coefficients_size();
  const auto& mplicand = request->multiplicand();
  const int mplicand_columns_size = mplicand.columns_size();
  const int mplicand_rows_size = mplicand.columns(0).coefficients_size();

  if (mplier_columns_size != mplicand_rows_size) {  // very shallow validation
    return Status(StatusCode::OUT_OF_RANGE, "1-factor cols != 2-factor rows");
  }

  MatrixXd mplier_eigen(mplier_rows_size, mplier_columns_size);
  for (int i = 0; i != mplier_columns_size; ++i) {
    for (int j = 0; j != mplier_rows_size; ++j) {
      mplier_eigen(j, i) = mplier.columns(i).coefficients(j);
    }
  }
  MatrixXd mplicand_eigen(mplicand_rows_size, mplicand_columns_size);
  for (int i = 0; i != mplicand_columns_size; ++i) {
    for (int j = 0; j != mplicand_rows_size; ++j) {
      mplicand_eigen(j, i) = mplicand.columns(i).coefficients(j);
    }
  }

  MatrixXd result_eigen = mplier_eigen * mplicand_eigen;

  auto response_result = response->mutable_result();
  for (Index ic = 0, ec = result_eigen.cols(); ec != ic; ++ic) {
    auto response_result_column = response_result->add_columns();
    for (Index ir = 0, er = result_eigen.rows(); er != ir; ++ir) {
      response_result_column->add_coefficients(result_eigen(ir, ic));
    }
  }

  return Status::OK;
}
