package alerting

import (
	"fmt"

	"github.com/grafana/grafana/pkg/services/provisioning/values"
)

type configVersion struct {
	APIVersion values.Int64Value `json:"apiVersion" yaml:"apiVersion"`
}

type OrgID int64

type AlertingFile struct {
	configVersion
	Filename            string
	Groups              []AlertRuleGroup
	DeleteRules         []RuleDelete
	ContactPoints       []ContactPoint
	DeleteContactPoints []DeleteContactPoint
	Policies            []NotificiationPolicy
	ResetPolicies       []OrgID
	MuteTimes           []MuteTime
	DeleteMuteTimes     []DeleteMuteTime
}

type AlertingFileV1 struct {
	configVersion
	Filename            string
	Groups              []AlertRuleGroupV1      `json:"groups" yaml:"groups"`
	DeleteRules         []RuleDeleteV1          `json:"deleteRules" yaml:"deleteRules"`
	ContactPoints       []ContactPointV1        `json:"contactPoints" yaml:"contactPoints"`
	DeleteContactPoints []DeleteContactPointV1  `json:"deleteContactPoints" yaml:"deleteContactPoints"`
	Policies            []NotificiationPolicyV1 `json:"policies" yaml:"policies"`
	ResetPolicies       []values.Int64Value     `json:"resetPolicies" yaml:"resetPolicies"`
	MuteTimes           []MuteTimeV1            `json:"muteTimes" yaml:"muteTimes"`
	DeleteMuteTimes     []DeleteMuteTimeV1      `json:"deleteMuteTimes" yaml:"deleteMuteTimes"`
}

func (fileV1 *AlertingFileV1) MapToModel() (AlertingFile, error) {
	alertingFile := AlertingFile{}
	alertingFile.Filename = fileV1.Filename
	if err := fileV1.mapRules(&alertingFile); err != nil {
		return AlertingFile{}, fmt.Errorf("failure parsing rules: %w", err)
	}
	if err := fileV1.mapContactPoint(&alertingFile); err != nil {
		return AlertingFile{}, fmt.Errorf("failure parsing contact points: %w", err)
	}
	fileV1.mapPolicies(&alertingFile)
	if err := fileV1.mapMuteTimes(&alertingFile); err != nil {
		return AlertingFile{}, fmt.Errorf("failure parsing mute times: %w", err)
	}
	return alertingFile, nil
}

func (fileV1 *AlertingFileV1) mapMuteTimes(alertingFile *AlertingFile) error {
	for _, mtV1 := range fileV1.MuteTimes {
		alertingFile.MuteTimes = append(alertingFile.MuteTimes, mtV1.mapToModel())
	}
	for _, deleteV1 := range fileV1.DeleteMuteTimes {
		delReq, err := deleteV1.mapToModel()
		if err != nil {
			return err
		}
		alertingFile.DeleteMuteTimes = append(alertingFile.DeleteMuteTimes, delReq)
	}
	return nil
}

func (fileV1 *AlertingFileV1) mapPolicies(alertingFile *AlertingFile) {
	for _, npV1 := range fileV1.Policies {
		alertingFile.Policies = append(alertingFile.Policies, npV1.mapToModel())
	}
	for _, orgIDV1 := range fileV1.ResetPolicies {
		alertingFile.ResetPolicies = append(alertingFile.ResetPolicies, OrgID(orgIDV1.Value()))
	}
}

func (fileV1 *AlertingFileV1) mapContactPoint(alertingFile *AlertingFile) error {
	for _, dcp := range fileV1.DeleteContactPoints {
		alertingFile.DeleteContactPoints = append(alertingFile.DeleteContactPoints, dcp.MapToModel())
	}
	for _, contactPointV1 := range fileV1.ContactPoints {
		contactPoint, err := contactPointV1.MapToModel()
		if err != nil {
			return err
		}
		alertingFile.ContactPoints = append(alertingFile.ContactPoints, contactPoint)
	}
	return nil
}

func (fileV1 *AlertingFileV1) mapRules(alertingFile *AlertingFile) error {
	for _, groupV1 := range fileV1.Groups {
		group, err := groupV1.MapToModel()
		if err != nil {
			return err
		}
		alertingFile.Groups = append(alertingFile.Groups, group)
	}
	for _, ruleDeleteV1 := range fileV1.DeleteRules {
		orgID := ruleDeleteV1.OrgID.Value()
		if orgID < 1 {
			orgID = 1
		}
		ruleDelete := RuleDelete{
			UID:   ruleDeleteV1.UID.Value(),
			OrgID: orgID,
		}
		alertingFile.DeleteRules = append(alertingFile.DeleteRules, ruleDelete)
	}
	return nil
}
