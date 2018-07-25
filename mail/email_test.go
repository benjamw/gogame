package mail

import (
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/benjamw/golibs/test"
	"github.com/mailgun/mailgun-go"

	"github.com/benjamw/gogame/config"
)

var (
	email = config.TestToEmail
)

func TestMain(m *testing.M) {
	test.InitCtx()

	config.GenRoot()
	config.RootURL = "UNIT_TESTING"

	runVal := m.Run()
	test.ReleaseCtx()
	os.Exit(runVal)
}

func TestParseTemplate(t *testing.T) {
	params := make(map[string]interface{}, 0)
	split, err := parseTemplates("welcome", params)
	if err != nil {
		t.Fatalf("mail.parseTemplates threw an error: %v", err)
	}
	if len(split) != 3 {
		t.Fatalf("mail.parseTemplates did not return the correct number of items. Wanted: 3; Got: %d", len(split))
	}
}

func TestFromTemplate(t *testing.T) {
	ctx := test.GetCtx()

	testMg := createMailgunTestMock(t)

	err := FromTemplate(ctx, "welcome", []string{email}, nil, testMg)
	if err != nil {
		t.Fatalf("mail.FromTemplate threw an error: %v", err)
	}
}

// MailgunMock is a mock implementation of Mailgun.

func createMailgunTestMock(t *testing.T) mailgun.Mailgun {

	// make and configure a mocked Mailgun
	mockedMailgun := &MailgunMock{
		AddBounceFunc: func(address string, code string, error string) error {
			panic("TODO: mock out the AddBounce method")
		},
		ApiBaseFunc: func() string {
			panic("TODO: mock out the ApiBase method")
		},
		ApiKeyFunc: func() string {
			panic("TODO: mock out the ApiKey method")
		},
		ChangeCredentialPasswordFunc: func(id string, password string) error {
			panic("TODO: mock out the ChangeCredentialPassword method")
		},
		ClientFunc: func() *http.Client {
			panic("TODO: mock out the Client method")
		},
		CreateCampaignFunc: func(name string, id string) error {
			panic("TODO: mock out the CreateCampaign method")
		},
		CreateComplaintFunc: func(in1 string) error {
			panic("TODO: mock out the CreateComplaint method")
		},
		CreateCredentialFunc: func(login string, password string) error {
			panic("TODO: mock out the CreateCredential method")
		},
		CreateDomainFunc: func(name string, smtpPassword string, spamAction string, wildcard bool) error {
			panic("TODO: mock out the CreateDomain method")
		},
		CreateListFunc: func(in1 mailgun.List) (mailgun.List, error) {
			panic("TODO: mock out the CreateList method")
		},
		CreateMemberFunc: func(merge bool, addr string, prototype mailgun.Member) error {
			panic("TODO: mock out the CreateMember method")
		},
		CreateMemberListFunc: func(subscribed *bool, addr string, newMembers []interface{}) error {
			panic("TODO: mock out the CreateMemberList method")
		},
		CreateRouteFunc: func(in1 mailgun.Route) (mailgun.Route, error) {
			panic("TODO: mock out the CreateRoute method")
		},
		CreateWebhookFunc: func(kind string, url string) error {
			panic("TODO: mock out the CreateWebhook method")
		},
		DeleteBounceFunc: func(address string) error {
			panic("TODO: mock out the DeleteBounce method")
		},
		DeleteCampaignFunc: func(id string) error {
			panic("TODO: mock out the DeleteCampaign method")
		},
		DeleteComplaintFunc: func(in1 string) error {
			panic("TODO: mock out the DeleteComplaint method")
		},
		DeleteCredentialFunc: func(id string) error {
			panic("TODO: mock out the DeleteCredential method")
		},
		DeleteDomainFunc: func(name string) error {
			panic("TODO: mock out the DeleteDomain method")
		},
		DeleteListFunc: func(in1 string) error {
			panic("TODO: mock out the DeleteList method")
		},
		DeleteMemberFunc: func(Member string, list string) error {
			panic("TODO: mock out the DeleteMember method")
		},
		DeleteRouteFunc: func(in1 string) error {
			panic("TODO: mock out the DeleteRoute method")
		},
		DeleteStoredMessageFunc: func(id string) error {
			panic("TODO: mock out the DeleteStoredMessage method")
		},
		DeleteTagFunc: func(tag string) error {
			panic("TODO: mock out the DeleteTag method")
		},
		DeleteWebhookFunc: func(kind string) error {
			panic("TODO: mock out the DeleteWebhook method")
		},
		DomainFunc: func() string {
			panic("TODO: mock out the Domain method")
		},
		GetBouncesFunc: func(limit int, skip int) (int, []mailgun.Bounce, error) {
			panic("TODO: mock out the GetBounces method")
		},
		GetCampaignsFunc: func() (int, []mailgun.Campaign, error) {
			panic("TODO: mock out the GetCampaigns method")
		},
		GetComplaintsFunc: func(limit int, skip int) (int, []mailgun.Complaint, error) {
			panic("TODO: mock out the GetComplaints method")
		},
		GetCredentialsFunc: func(limit int, skip int) (int, []mailgun.Credential, error) {
			panic("TODO: mock out the GetCredentials method")
		},
		GetDomainsFunc: func(limit int, skip int) (int, []mailgun.Domain, error) {
			panic("TODO: mock out the GetDomains method")
		},
		GetListByAddressFunc: func(in1 string) (mailgun.List, error) {
			panic("TODO: mock out the GetListByAddress method")
		},
		GetListsFunc: func(limit int, skip int, filter string) (int, []mailgun.List, error) {
			panic("TODO: mock out the GetLists method")
		},
		GetMemberByAddressFunc: func(MemberAddr string, listAddr string) (mailgun.Member, error) {
			panic("TODO: mock out the GetMemberByAddress method")
		},
		GetMembersFunc: func(limit int, skip int, subfilter *bool, address string) (int, []mailgun.Member, error) {
			panic("TODO: mock out the GetMembers method")
		},
		GetRouteByIDFunc: func(in1 string) (mailgun.Route, error) {
			panic("TODO: mock out the GetRouteByID method")
		},
		GetRoutesFunc: func(limit int, skip int) (int, []mailgun.Route, error) {
			panic("TODO: mock out the GetRoutes method")
		},
		GetSingleBounceFunc: func(address string) (mailgun.Bounce, error) {
			panic("TODO: mock out the GetSingleBounce method")
		},
		GetSingleComplaintFunc: func(address string) (mailgun.Complaint, error) {
			panic("TODO: mock out the GetSingleComplaint method")
		},
		GetSingleDomainFunc: func(domain string) (mailgun.Domain, []mailgun.DNSRecord, []mailgun.DNSRecord, error) {
			panic("TODO: mock out the GetSingleDomain method")
		},
		GetStatsFunc: func(limit int, skip int, startDate *time.Time, event ...string) (int, []mailgun.Stat, error) {
			panic("TODO: mock out the GetStats method")
		},
		GetStoredMessageFunc: func(id string) (mailgun.StoredMessage, error) {
			panic("TODO: mock out the GetStoredMessage method")
		},
		GetStoredMessageForURLFunc: func(url string) (mailgun.StoredMessage, error) {
			panic("TODO: mock out the GetStoredMessageForURL method")
		},
		GetStoredMessageRawFunc: func(id string) (mailgun.StoredMessageRaw, error) {
			panic("TODO: mock out the GetStoredMessageRaw method")
		},
		GetStoredMessageRawForURLFunc: func(url string) (mailgun.StoredMessageRaw, error) {
			panic("TODO: mock out the GetStoredMessageRawForURL method")
		},
		GetTagFunc: func(tag string) (mailgun.TagItem, error) {
			panic("TODO: mock out the GetTag method")
		},
		GetUnsubscribesFunc: func(limit int, skip int) (int, []mailgun.Unsubscription, error) {
			panic("TODO: mock out the GetUnsubscribes method")
		},
		GetUnsubscribesByAddressFunc: func(in1 string) (int, []mailgun.Unsubscription, error) {
			panic("TODO: mock out the GetUnsubscribesByAddress method")
		},
		GetWebhookByTypeFunc: func(kind string) (string, error) {
			panic("TODO: mock out the GetWebhookByType method")
		},
		GetWebhooksFunc: func() (map[string]string, error) {
			panic("TODO: mock out the GetWebhooks method")
		},
		ListEventsFunc: func(in1 *mailgun.EventsOptions) *mailgun.EventIterator {
			panic("TODO: mock out the ListEvents method")
		},
		ListTagsFunc: func(in1 *mailgun.TagOptions) *mailgun.TagIterator {
			panic("TODO: mock out the ListTags method")
		},
		NewEventIteratorFunc: func() *mailgun.EventIterator {
			panic("TODO: mock out the NewEventIterator method")
		},
		NewMIMEMessageFunc: func(body io.ReadCloser, to ...string) *mailgun.Message {
			panic("TODO: mock out the NewMIMEMessage method")
		},
		NewMessageFunc: func(from string, subject string, text string, to ...string) *mailgun.Message {
			panic("TODO: mock out the NewMessage method")
		},
		ParseAddressesFunc: func(addresses ...string) ([]string, []string, error) {
			panic("TODO: mock out the ParseAddresses method")
		},
		PollEventsFunc: func(in1 *mailgun.EventsOptions) *mailgun.EventPoller {
			panic("TODO: mock out the PollEvents method")
		},
		PublicApiKeyFunc: func() string {
			panic("TODO: mock out the PublicApiKey method")
		},
		RemoveUnsubscribeFunc: func(in1 string) error {
			panic("TODO: mock out the RemoveUnsubscribe method")
		},
		RemoveUnsubscribeWithTagFunc: func(a string, t string) error {
			panic("TODO: mock out the RemoveUnsubscribeWithTag method")
		},
		SendFunc: func(m *mailgun.Message) (string, string, error) {
			panic("TODO: mock out the Send method")
		},
		SetAPIBaseFunc: func(url string) {
			panic("TODO: mock out the SetAPIBase method")
		},
		SetClientFunc: func(client *http.Client) {
			panic("TODO: mock out the SetClient method")
		},
		UnsubscribeFunc: func(address string, tag string) error {
			panic("TODO: mock out the Unsubscribe method")
		},
		UpdateCampaignFunc: func(oldId string, name string, newId string) error {
			panic("TODO: mock out the UpdateCampaign method")
		},
		UpdateListFunc: func(in1 string, in2 mailgun.List) (mailgun.List, error) {
			panic("TODO: mock out the UpdateList method")
		},
		UpdateMemberFunc: func(Member string, List string, prototype mailgun.Member) (mailgun.Member, error) {
			panic("TODO: mock out the UpdateMember method")
		},
		UpdateRouteFunc: func(in1 string, in2 mailgun.Route) (mailgun.Route, error) {
			panic("TODO: mock out the UpdateRoute method")
		},
		UpdateWebhookFunc: func(kind string, url string) error {
			panic("TODO: mock out the UpdateWebhook method")
		},
		ValidateEmailFunc: func(email string) (mailgun.EmailVerification, error) {
			panic("TODO: mock out the ValidateEmail method")
		},
		VerifyWebhookRequestFunc: func(req *http.Request) (bool, error) {
			panic("TODO: mock out the VerifyWebhookRequest method")
		},
	}

	return mockedMailgun
}
