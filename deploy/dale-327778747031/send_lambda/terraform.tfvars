terragrunt = {
  terraform {
    source = "../../src//send_lambda"
  }

  remote_state {
    backend = "s3"
    config {
      bucket = "prm-327778747031-terraform-states"
      key = "dale/send_lambda/terraform.tfstate"
      region = "eu-west-2"
      encrypt = true
    }
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# MODULE PARAMETERS
# These are the variables we have to pass in to use the module specified in the terragrunt configuration above
# ---------------------------------------------------------------------------------------------------------------------

aws_region = "eu-west-2"
environment = "dale"
# Path has 5 extra .. components in as Terragrunt copies the modules into a cache in several sub-directories of this directory
lambda_zip = "../../../../../../../lambda/send/lambda.zip"