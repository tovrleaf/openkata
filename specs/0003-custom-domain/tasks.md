# Tasks: Custom Domain with CloudFront

## Tasks

### 1. Create Route 53 hosted zone script
- **Status**: Pending
- **Goal**: Script that creates the hosted zone and outputs
  nameservers for the user to configure at registrar
- **Boundary**: `infra/`
- **Depends**: None
- **Done when**: Script creates hosted zone and prints NS
  records

### 2. Create ACM certificate script
- **Status**: Pending
- **Goal**: Script that requests ACM certificate for
  openkata.dev and *.openkata.dev, adds DNS validation
  records to Route 53, and waits for validation
- **Boundary**: `infra/`
- **Depends**: 1 (hosted zone must exist, nameservers must
  be pointed — user action between tasks)
- **Done when**: Script requests cert and adds CNAME
  validation records

### 3. Create CloudFront distribution script
- **Status**: Pending
- **Goal**: Script that creates CloudFront distribution
  with the Lambda Function URL as origin, attaches the ACM
  certificate, and configures both openkata.dev and
  www.openkata.dev as aliases
- **Boundary**: `infra/`
- **Depends**: 2 (certificate must be validated)
- **Done when**: Script creates distribution and outputs
  the CloudFront domain name

### 4. Create Route 53 alias records script
- **Status**: Pending
- **Goal**: Script that creates A/AAAA alias records
  pointing openkata.dev and www.openkata.dev to the
  CloudFront distribution
- **Boundary**: `infra/`
- **Depends**: 3
- **Done when**: Script creates DNS records

### 5. Update IAM admin policy
- **Status**: Pending
- **Goal**: Add Route 53, ACM, and CloudFront permissions
  to iam-admin-policy.json
- **Boundary**: `infra/iam-admin-policy.json`
- **Depends**: None
- **Done when**: Policy covers all actions used by the
  scripts

## Progress Log
