/*******************************************************************************
 * Copyright 2019 Dell Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/

package certificates

import (
	"fmt"

	"github.com/edgexfoundry/edgex-go/internal/security/secrets/seed"

	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
)

type CertificateType int

const (
	RootCertificate CertificateType = 1
	TLSCertificate  CertificateType = 2
)

type CertificateGenerator interface {
	Generate() error
}

func NewCertificateGenerator(t CertificateType, certificateSeed seed.CertificateSeed, w FileWriter, logger logger.LoggingClient) (CertificateGenerator, error) {
	switch t {
	case RootCertificate:
		return rootCertGenerator{certificateSeed: certificateSeed, writer: w, logger: logger}, nil
	case TLSCertificate:
		return tlsCertGenerator{certificateSeed: certificateSeed, writer: w, logger: logger}, nil
	default:
		return nil, fmt.Errorf("unknown CertificateType %v", t)
	}
}
