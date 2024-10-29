//
// Copyright 2024 The Sigstore Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/sigstore/rekor-monitor/pkg/identity"
	"github.com/sigstore/rekor-monitor/pkg/notifications"
)

type IdentityMonitorConfiguration struct {
	StartIndex                *int                                     `yaml:"startIndex"`
	EndIndex                  *int                                     `yaml:"endIndex"`
	MonitoredValues           identity.MonitoredValues                 `yaml:"monitoredValues"`
	ServerURL                 string                                   `yaml:"serverURL"`
	OutputIdentitiesFile      string                                   `yaml:"outputIdentities"`
	LogInfoFile               string                                   `yaml:"logInfoFile"`
	IdentityMetadataFile      *string                                  `yaml:"identityMetadataFile"`
	GitHubIssue               *notifications.GitHubIssueInput          `yaml:"githubIssue"`
	EmailNotificationSMTP     *notifications.EmailNotificationInput    `yaml:"emailNotificationSMTP"`
	EmailNotificationMailgun  *notifications.MailgunNotificationInput  `yaml:"emailNotificationMailgun"`
	EmailNotificationSendGrid *notifications.SendGridNotificationInput `yaml:"emailNotificationSendGrid"`
	Interval                  *time.Duration                           `yaml:"interval"`
}

func (config IdentityMonitorConfiguration) TriggerNotifications(identities []identity.MonitoredIdentity) error {
	// update this as new notification platforms are implemented within rekor-monitor
	if config.GitHubIssue != nil {
		err := config.GitHubIssue.Send(context.Background(), identities)
		if err != nil {
			return fmt.Errorf("error creating new GitHub issue: %v", err)
		}
	}

	if config.EmailNotificationSMTP != nil {
		err := config.EmailNotificationSMTP.Send(context.Background(), identities)
		if err != nil {
			return fmt.Errorf("error sending email via SMTP: %v", err)
		}
	}

	if config.EmailNotificationSendGrid != nil {
		err := config.EmailNotificationSMTP.Send(context.Background(), identities)
		if err != nil {
			return fmt.Errorf("error sending email via SendGrid: %v", err)
		}
	}

	if config.EmailNotificationMailgun != nil {
		err := config.EmailNotificationSMTP.Send(context.Background(), identities)
		if err != nil {
			return fmt.Errorf("error sending email via SendGrid: %v", err)
		}
	}

	return nil
}
