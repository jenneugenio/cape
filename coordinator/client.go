package coordinator

import (
	"context"

	"github.com/capeprivacy/cape/models"
	"github.com/manifoldco/go-base64"

	errors "github.com/capeprivacy/cape/partyerrors"

	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/coordinator/graph/model"
)

// NetworkCause occurs when the client cannot reach the server
var NetworkCause = errors.NewCause(errors.RequestTimeoutCategory, "network_error")

// Client is a wrapper around the graphql client that
// connects to the coordinator and sends queries
type Client struct {
	transport ClientTransport
}

// NewClient returns a new client that connects to the given
// the configured transport
func NewClient(transport ClientTransport) *Client {
	return &Client{
		transport: transport,
	}
}

type MeResponse struct {
	User *models.User `json:"me"`
}

// Me returns the user of the current authenticated session
func (c *Client) Me(ctx context.Context) (*models.User, error) {
	var resp MeResponse

	err := c.transport.Raw(ctx, `
		query Me {
			me {
				id
				email
				name
			}
		}
	`, nil, &resp)
	if err != nil {
		return nil, err
	}

	return resp.User, nil
}

// UserResponse is a User with an extra Roles field that maps to the
// GraphQL type.
type UserResponse struct {
	*models.User
	Role models.Role `json:"role"`
}

// GetUser returns a user and it's roles!
func (c *Client) GetUser(ctx context.Context, id string) (*UserResponse, error) {
	var resp struct {
		User UserResponse `json:"user"`
	}

	variables := make(map[string]interface{})
	variables["id"] = id

	err := c.transport.Raw(ctx, `
		query User($id: String!) {
			user(id: $id) {
				id
				name
				email
				role {
					id
					label
				}
			}
		}
	`, variables, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.User, nil
}

// CreateUser creates a user and returns it
func (c *Client) CreateUser(ctx context.Context, name models.Name, email models.Email) (*models.User, models.Password, error) {
	var resp struct {
		Response struct {
			Password models.Password `json:"password"`
			User     *models.User    `json:"user"`
		} `json:"createUser"`
	}

	variables := make(map[string]interface{})
	variables["name"] = name
	variables["email"] = email

	err := c.transport.Raw(ctx, `
		mutation CreateUser ($name: Name!, $email: ModelEmail!) {
			createUser(input: { name: $name, email: $email }) {
				password
				user {
					id
					name
					email
				}
			}
		}
	`, variables, &resp)

	if err != nil {
		return nil, "", err
	}

	return resp.Response.User, resp.Response.Password, nil
}

// ListUsers returns all of the users in the database
func (c *Client) ListUsers(ctx context.Context) ([]*models.User, error) {
	var resp struct {
		Users []*models.User `json:"users"`
	}

	err := c.transport.Raw(ctx, `
		query Users {
			users {
				id
				name
				email
				role {
					id
					label
				}
			}
		}
	`, nil, &resp)

	if err != nil {
		return nil, err
	}

	return resp.Users, nil
}

func (c *Client) Authenticated() bool {
	return c.transport.Authenticated()
}

// EmailLogin calls the CreateLoginSession and CreateAuthSession mutations
func (c *Client) EmailLogin(ctx context.Context, email models.Email, password models.Password) (*models.Session, error) {
	return c.transport.EmailLogin(ctx, email, password)
}

func (c *Client) TokenLogin(ctx context.Context, token *auth.APIToken) (*models.Session, error) {
	return c.transport.TokenLogin(ctx, token)
}

// Logout calls the deleteSession mutation
func (c *Client) Logout(ctx context.Context, authToken *base64.Value) error {
	return c.transport.Logout(ctx, authToken)
}

// SessionToken returns the client's current session token
func (c *Client) SessionToken() *base64.Value {
	return c.transport.Token()
}

// Role Routes
type MyRoleResponse struct {
	Role models.Role `json:"myRole"`
}

func (c *Client) MyRole(ctx context.Context) (*models.Role, error) {
	var resp MyRoleResponse

	err := c.transport.Raw(ctx, `
		query MyRole() {
			myRole() {
				id,
				label
			}
		}
	`, nil, &resp)

	if err != nil {
		return nil, err
	}

	return &resp.Role, nil
}

func (c *Client) MyProjectRole(ctx context.Context, project models.Label) (*models.Role, error) {
	var resp struct {
		Role models.Role `json:"myRole"`
	}

	variables := make(map[string]interface{})
	variables["project_label"] = project

	err := c.transport.Raw(ctx, `
		query MyRole($project_label: ModelLabel) {
			myRole(project_label: $project_label) {
				id,
				label
			}
		}
	`, variables, &resp)

	if err != nil {
		return nil, err
	}

	return &resp.Role, nil
}

func (c *Client) SetOrgRole(ctx context.Context, user models.Email, role models.Label) error {
	variables := make(map[string]interface{})
	variables["user_email"] = user
	variables["role_label"] = role

	err := c.transport.Raw(ctx, `
		mutation SetOrgRole($user_email: ModelEmail!, $role_label: ModelLabel!) {
			setOrgRole(user_email: $user_email, role_label: $role_label) { id }
		}
	`, variables, nil)

	return err
}

func (c *Client) SetProjectRole(ctx context.Context, user models.Email, project models.Label, role models.Label) error {
	variables := make(map[string]interface{})
	variables["user_email"] = user
	variables["role_label"] = role
	variables["project_label"] = project

	err := c.transport.Raw(ctx, `
		mutation SetProjectRole($user_email: ModelEmail!, $project_label: ModelLabel!, $role_label: ModelLabel!) {
			setProjectRole(user_email: $user_email, project_label: $project_label, role_label: $role_label) { id }
		}
	`, variables, nil)

	return err
}

// GetRoleByLabel returns a specific role by label
func (c *Client) GetRoleByLabel(ctx context.Context, label models.Label) (*models.Role, error) {
	var resp struct {
		Role models.Role `json:"roleByLabel"`
	}

	variables := make(map[string]interface{})
	variables["label"] = label

	err := c.transport.Raw(ctx, `
		query RoleByLabel($label: Label!) {
			roleByLabel(label: $label) {
				id
				label
			}
		}
	`, variables, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Role, nil
}

// GetUsers returns all users for the given emails
func (c *Client) GetUsers(ctx context.Context, emails []models.Email) ([]*models.User, error) {
	var resp struct {
		Users []*models.User `json:"identities"`
	}

	variables := make(map[string]interface{})
	variables["emails"] = emails

	err := c.transport.Raw(ctx, `
		query GetUsers($emails: [Email!]) {
			identities(emails: $emails) {
				id
				email
			}
		}
	`, variables, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Users, nil
}

type CreateTokenMutation struct {
	Secret models.Password `json:"secret"`
	Token  *models.Token   `json:"token"`
}

type CreateTokenResponse struct {
	Response *CreateTokenMutation `json:"createToken"`
}

// CreateToken creates a new API token for the provided user. You can pass nil and it will return a token for you
func (c *Client) CreateToken(ctx context.Context, user *models.User) (*auth.APIToken, *models.Token, error) {
	// If the user provides no user, we will make a token for the current session user
	if user == nil {
		i, err := c.Me(ctx)
		if err != nil {
			return nil, nil, err
		}

		user = i
	}

	variables := make(map[string]interface{})
	variables["user_id"] = user.ID

	resp := &CreateTokenResponse{}
	err := c.transport.Raw(ctx, `
		mutation CreateToken($user_id: String!) {
			createToken(input: { user_id: $user_id }) {
				secret
				token {
					id
				}
			}
        }
    `, variables, resp)
	if err != nil {
		return nil, nil, err
	}

	secret, err := auth.FromPassword(resp.Response.Secret)
	if err != nil {
		return nil, nil, err
	}

	token, err := auth.NewAPIToken(secret, resp.Response.Token.ID)
	return token, resp.Response.Token, err
}

type ListTokensResponse struct {
	IDs []string `json:"tokens"`
}

// ListTokens lists all of the auth tokens for the provided user
func (c *Client) ListTokens(ctx context.Context, user *models.User) ([]string, error) {
	// If the user provides no user, we will make a token for the current session user
	if user == nil {
		i, err := c.Me(ctx)
		if err != nil {
			return nil, err
		}

		user = i
	}

	var resp ListTokensResponse

	variables := make(map[string]interface{})
	variables["user_id"] = user.ID

	err := c.transport.Raw(ctx, `
		query Tokens($user_id: String!) {
			tokens(user_id: $user_id)
		}
    `, variables, &resp)

	if err != nil {
		return nil, err
	}

	return resp.IDs, nil
}

// RemoveToken removes the provided token from the database
func (c *Client) RemoveToken(ctx context.Context, tokenID string) error {
	variables := make(map[string]interface{})
	variables["id"] = tokenID

	return c.transport.Raw(ctx, `
		mutation RemoveToken($id: String!) {
			removeToken(id: $id)
		}
    `, variables, nil)
}

type CreateProjectResponse struct {
	Project *models.Project `json:"createProject"`
}

func (c *Client) CreateProject(
	ctx context.Context,
	name models.ProjectDisplayName,
	label *models.Label,
	desc models.ProjectDescription) (*models.Project, error) {
	createReq := &model.CreateProjectRequest{
		Name:        name,
		Label:       label,
		Description: desc,
	}

	variables := make(map[string]interface{})
	variables["create_project"] = createReq

	var resp CreateProjectResponse

	err := c.transport.Raw(ctx, `
		mutation CreateProject($create_project: CreateProjectRequest!) {
			createProject(project: $create_project) {
				id,
				name,
				label,
				description,
				status
			}
		}
	`, variables, &resp)

	if err != nil {
		return nil, err
	}

	return resp.Project, nil
}

type ListProjectsResponse struct {
	Projects []*models.Project `json:"projects"`
}

func (c *Client) ListProjects(ctx context.Context, status models.ProjectStatus) ([]*models.Project, error) {
	variables := make(map[string]interface{})
	variables["status"] = status

	var resp ListProjectsResponse
	err := c.transport.Raw(ctx, `
		query ListProjects($status: ProjectStatus!) {
			projects(status: $status) {
				id,
				name,
				label,
				description,
				status
			}
		}
	`, variables, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Projects, nil
}

type GetProject struct {
	*models.Project
	Policy       *models.Policy   `json:"current_spec"`
	Contributors []GQLContributor `json:"contributors"`
}

type GetProjectResponse struct {
	GetProject GetProject `json:"project"`
}

func (c *Client) GetProject(ctx context.Context, id string, label *models.Label) (*GetProject, error) {
	var resp GetProjectResponse
	variables := make(map[string]interface{})
	if id != "" {
		variables["id"] = id
	}

	if label != nil {
		variables["label"] = label
	}

	err := c.transport.Raw(ctx, `
		query GetProjects($id: String, $label: ModelLabel) {
			project(id: $id, label: $label) {
				id,
				name,
				label,
				description,
				status,
				created_at,
				updated_at,

				current_spec {
					transformations
					rules
				}
					

				contributors {
					id
					user {
						id
						name
						email
					}
					role {
						id
						label
					}
				}
			}
		}
	`, variables, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.GetProject, nil
}

type UpdateProjectSpecResponseBody struct {
	*models.Project
	ProjectSpec *models.Policy `json:"current_spec"`
}

type UpdateProjectSpecResponse struct {
	UpdateProjectSpecResponseBody `json:"updateProjectSpec"`
}

func (c *Client) UpdateProjectSpec(ctx context.Context, projectLabel models.Label, spec *models.PolicyFile) (*models.Project, *models.Policy, error) {
	variables := make(map[string]interface{})
	variables["project"] = &projectLabel
	variables["projectSpecFile"] = spec

	var resp UpdateProjectSpecResponse
	err := c.transport.Raw(ctx, `
		mutation UpdateProjectSpec($project: ModelLabel, $projectSpecFile: ProjectSpecFile!) {
			updateProjectSpec(label: $project, request: $projectSpecFile) {
				name,
				label,
				description,
				status,
				current_spec {
					id,
					rules
				}
			}
		}
	`, variables, &resp)
	if err != nil {
		return nil, nil, err
	}

	p := resp.UpdateProjectSpecResponseBody.Project
	s := resp.UpdateProjectSpecResponseBody.ProjectSpec

	p.CurrentSpecID = s.ID
	return p, s, nil
}

type SuggestPolicyResponse struct {
	Suggestion models.Suggestion `json:"suggestProjectPolicy"`
}

func (c *Client) SuggestPolicy(
	ctx context.Context,
	projectLabel models.Label,
	name models.ProjectDisplayName,
	description models.ProjectDescription,
	spec *models.PolicyFile) (*models.Suggestion, error) {
	variables := make(map[string]interface{})
	variables["project"] = &projectLabel
	variables["projectSpecFile"] = spec
	variables["name"] = name
	variables["description"] = description

	var resp SuggestPolicyResponse

	err := c.transport.Raw(ctx, `
		mutation SuggestProjectPolicy(
			$project: ModelLabel!, 
			$projectSpecFile: ProjectSpecFile!,
			$name: String!
			$description: String!) {
			suggestProjectPolicy(label: $project, name: $name, description: $description, request: $projectSpecFile) {
				id
				state
				title
				description
				created_at
				updated_at
			}
		}
	`, variables, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Suggestion, nil
}

type GetProjectSuggestionsResponse struct {
	Suggestions []models.Suggestion `json:"getProjectSuggestions"`
}

func (c *Client) GetProjectSuggestions(ctx context.Context, projectLabel models.Label) ([]models.Suggestion, error) {
	variables := make(map[string]interface{})
	variables["project"] = &projectLabel

	var resp GetProjectSuggestionsResponse

	err := c.transport.Raw(ctx, `
		mutation GetProjectSuggestions($project: ModelLabel!) {
			getProjectSuggestions(label: $project) {
				state
				title
				id
				created_at
				updated_at
			}
		}
	`, variables, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Suggestions, nil
}

type ProjectSuggestion struct {
	*models.Suggestion
	Policy  models.Policy  `json:"policy"`
	Project models.Project `json:"project"`
}

type GetProjectSuggestionResponse struct {
	SuggestionResponse ProjectSuggestion `json:"getProjectSuggestion"`
}

func (c *Client) GetProjectSuggestion(ctx context.Context, id string) (*ProjectSuggestion, error) {
	variables := make(map[string]interface{})
	variables["id"] = id

	var resp GetProjectSuggestionResponse

	err := c.transport.Raw(ctx, `
		mutation GetProjectSuggestion($id: String!) {
			getProjectSuggestion(id: $id) {
				state
				title
				description
				id
				project {
					id
					name
					label
					description
					status
				}
				policy {
					id
					rules
					transformations
				}
				created_at
				updated_at
			}
		}
	`, variables, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.SuggestionResponse, nil
}

type RejectSuggestionResponse struct {
	Project models.Project `json:"rejectProjectSuggestion"`
}

func (c *Client) RejectSuggestion(ctx context.Context, suggestion models.Suggestion) error {
	variables := make(map[string]interface{})
	variables["id"] = suggestion.ID

	var resp RejectSuggestionResponse

	err := c.transport.Raw(ctx, `
		mutation RejectProjectSuggestion($id: String!) {
			rejectProjectSuggestion(id: $id) {
				id
			}
		}
	`, variables, &resp)
	if err != nil {
		return err
	}

	return nil
}

type ApproveSuggestionResponse struct {
	Project models.Project `json:"approveProjectSuggestion"`
}

func (c *Client) ApproveSuggestion(ctx context.Context, suggestion models.Suggestion) error {
	variables := make(map[string]interface{})
	variables["id"] = suggestion.ID

	var resp ApproveSuggestionResponse

	err := c.transport.Raw(ctx, `
		mutation ApproveProjectSuggestion($id: String!) {
			approveProjectSuggestion(id: $id) {
				id
			}
		}
	`, variables, &resp)
	if err != nil {
		return err
	}

	return nil
}

type UpdateProjectResponse struct {
	Project *models.Project `json:"updateProject"`
}

func (c *Client) UpdateProject(
	ctx context.Context,
	id string,
	label *models.Label,
	name *models.ProjectDisplayName,
	desc *models.ProjectDescription) (*models.Project, error) {
	updateReq := &model.UpdateProjectRequest{
		Name:        name,
		Description: desc,
	}

	variables := map[string]interface{}{
		"id":             id,
		"label":          label,
		"update_project": updateReq,
	}

	var resp UpdateProjectResponse

	err := c.transport.Raw(ctx, `
		mutation UpdateProject($id: String, $label: ModelLabel, $update_project: UpdateProjectRequest!) {
			updateProject(id: $id, label: $label, update: $update_project) {
				id,
				name,
				label,
				description,
				status
			}
		}
	`, variables, &resp)

	if err != nil {
		return nil, err
	}

	return resp.Project, nil
}

type UpdateContributorResponse struct {
	*models.Contributor `json:"updateContributor"`
	User                *models.User `json:"user"`
}

func (c *Client) AddContributor(ctx context.Context, project models.Project, user models.Email, role models.Label) (*models.Contributor, error) {
	variables := map[string]interface{}{
		"project_label": project.Label,
		"email":         user,
		"role":          role,
	}

	var resp UpdateContributorResponse

	query := `mutation updateContributor($project_label: ModelLabel!, $email: ModelEmail!, $role: ModelLabel!) {
		updateContributor(project_label: $project_label, user_email: $email, role_label: $role) {
			id
			user {
			  id
			}
			project {
              id
            }
		}
	}`

	err := c.transport.Raw(ctx, query, variables, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Contributor, nil
}

type RemoveContributorResponse struct {
	Contributor models.Contributor `json:"removeContributor"`
}

func (c *Client) RemoveContributor(ctx context.Context, user models.User, project models.Project) (*models.Contributor, error) {
	variables := map[string]interface{}{
		"project_label": project.Label,
		"email":         user.Email,
	}

	var resp RemoveContributorResponse

	query := `mutation removeContributor($project_label: ModelLabel!, $email: ModelEmail!) {
		removeContributor(project_label: $project_label, user_email: $email) {
			id
		}
	}`

	err := c.transport.Raw(ctx, query, variables, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Contributor, nil
}

type GQLContributor struct {
	*models.Contributor
	User    models.User    `json:"user"`
	Project models.Project `json:"project"`
	Role    models.Role    `json:"role"`
}

type ListContributorsResponse struct {
	Contributors []GQLContributor `json:"listContributors"`
}

func (c *Client) ListContributors(ctx context.Context, project models.Project) ([]GQLContributor, error) {
	variables := map[string]interface{}{
		"project_label": project.Label,
	}

	var resp ListContributorsResponse
	query := `query listContributors($project_label: ModelLabel!) {
		listContributors(project_label: $project_label) {
			id,
			created_at,
			updated_at

			user {
				id
				name
				email
			}

			project {
				id
				label
			}

			role {
				id
				label
				system
			}
		}
	}`

	err := c.transport.Raw(ctx, query, variables, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Contributors, nil
}

func (c *Client) CreateRecovery(ctx context.Context, email models.Email) error {
	variables := map[string]interface{}{
		"email": email,
	}

	query := `mutation createRecovery($email: ModelEmail!) {
		createRecovery(input: { email: $email })
	}`

	return c.transport.Raw(ctx, query, variables, nil)
}

func (c *Client) AttemptRecovery(ctx context.Context, ID string, secret models.Password, newPassword models.Password) error {
	variables := map[string]interface{}{
		"new_password": newPassword.String(),
		"secret":       secret.String(),
		"id":           ID,
	}

	query := `mutation attemptRecovery($new_password: Password!, $secret: Password!, $id: String!) {
		attemptRecovery(input: {
			new_password: $new_password,
			secret: $secret,
			id: $id,
		})
	}`

	return c.transport.Raw(ctx, query, variables, nil)
}

type ListRecoveriesResponse struct {
	Recoveries []models.Recovery `json:"recoveries"`
}

func (c *Client) Recoveries(ctx context.Context) ([]models.Recovery, error) {
	var resp ListRecoveriesResponse
	query := `
		query ListRecoveries() {
			recoveries() {
				id,
				created_at,
				updated_at
			}
		}
	`
	err := c.transport.Raw(ctx, query, nil, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Recoveries, nil
}

func (c *Client) DeleteRecoveries(ctx context.Context, ids []string) error {
	variables := map[string]interface{}{
		"ids": ids,
	}

	query := `
		mutation DeleteRecoveries($ids: [String!]!) {
			deleteRecoveries(input: { ids: $ids })
		}
	`

	return c.transport.Raw(ctx, query, variables, nil)
}
