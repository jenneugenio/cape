package coordinator

import (
	"context"

	"github.com/capeprivacy/cape/models"
	"github.com/manifoldco/go-base64"

	errors "github.com/capeprivacy/cape/partyerrors"

	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/graph/model"
	"github.com/capeprivacy/cape/primitives"
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

// UserResponse is a primitive User with an extra Roles field that maps to the
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
func (c *Client) CreateUser(ctx context.Context, name models.Name, email models.Email) (*models.User, primitives.Password, error) {
	var resp struct {
		Response struct {
			Password primitives.Password `json:"password"`
			User     *models.User        `json:"user"`
		} `json:"createUser"`
	}

	variables := make(map[string]interface{})
	variables["name"] = name
	variables["email"] = email

	err := c.transport.Raw(ctx, `
		mutation CreateUser ($name: ModelName!, $email: ModelEmail!) {
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
		return nil, primitives.Password(""), err
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
func (c *Client) EmailLogin(ctx context.Context, email models.Email, password primitives.Password) (*primitives.Session, error) {
	return c.transport.EmailLogin(ctx, email, password)
}

func (c *Client) TokenLogin(ctx context.Context, token *auth.APIToken) (*primitives.Session, error) {
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

func (c *Client) MyRole(ctx context.Context) (*models.Role, error) {
	var resp struct {
		Role models.Role `json:"myRole"`
	}

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

// CreateRole creates a new role with a label
//func (c *Client) CreateRole(ctx context.Context, label primitives.Label, userIDs []string) (*primitives.Role, error) {
//	var resp struct {
//		Role primitives.Role `json:"createRole"`
//	}
//
//	variables := make(map[string]interface{})
//	variables["ids"] = userIDs
//	variables["label"] = label
//
//	err := c.transport.Raw(ctx, `
//		mutation CreateRole($label: Label!, $ids: [String!]) {
//			createRole(input: { label: $label, user_ids: $ids }) {
//				id
//				label
//				system
//			}
//		}
//	`, variables, &resp)
//	if err != nil {
//		return nil, err
//	}
//
//	return &resp.Role, nil
//}

// DeleteRole deletes a role with the given id
//func (c *Client) DeleteRole(ctx context.Context, id database.ID) error {
//	variables := make(map[string]interface{})
//	variables["id"] = id
//
//	return c.transport.Raw(ctx, `
//		mutation DeleteRole($id: ID!) {
//			deleteRole(input: { id: $id })
//		}
//	`, variables, nil)
//}

// GetRole returns a specific role
//func (c *Client) GetRole(ctx context.Context, id database.ID) (*primitives.Role, error) {
//	var resp struct {
//		Role primitives.Role `json:"role"`
//	}
//
//	variables := make(map[string]interface{})
//	variables["id"] = id
//
//	err := c.transport.Raw(ctx, `
//		query Role($id: ID!) {
//			role(id: $id) {
//				id
//				label
//			}
//		}
//	`, variables, &resp)
//	if err != nil {
//		return nil, err
//	}
//
//	return &resp.Role, nil
//}

// GetRoleByLabel returns a specific role by label
func (c *Client) GetRoleByLabel(ctx context.Context, label primitives.Label) (*models.Role, error) {
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

// GetMembersRole returns the members of a role
//func (c *Client) GetMembersRole(ctx context.Context, roleID database.ID) ([]*models.User, error) {
//	var resp struct {
//		Users []*models.User `json:"roleMembers"`
//	}
//
//	variables := make(map[string]interface{})
//	variables["role_id"] = roleID
//
//	err := c.transport.Raw(ctx, `
//		query GetMembersRole($role_id: ID!) {
//			roleMembers(role_id: $role_id) {
//				id
//				email
//			}
//		}
//	`, variables, &resp)
//	if err != nil {
//		return nil, err
//	}
//
//	return resp.Users, nil
//}

// AssignRole assigns a role to an user
//func (c *Client) AssignRole(ctx context.Context, userID string, roleID database.ID) (*model.Assignment, error) {
//	var resp struct {
//		Assignment model.Assignment `json:"assignRole"`
//	}
//
//	variables := make(map[string]interface{})
//	variables["role_id"] = roleID
//	variables["user_id"] = userID
//
//	err := c.transport.Raw(ctx, `
//		mutation AssignRole($role_id: ID!, $user_id: String!) {
//			assignRole(input: { role_id: $role_id, user_id: $user_id }) {
//				role {
//					id
//					label
//				}
//				user {
//					id
//					email
//				}
//			}
//		}
//	`, variables, &resp)
//	if err != nil {
//		return nil, err
//	}
//
//	return &resp.Assignment, nil
//}

// UnassignRole unassigns a role from an identity
//func (c *Client) UnassignRole(ctx context.Context, userID string, roleID database.ID) error {
//	variables := make(map[string]interface{})
//	variables["role_id"] = roleID
//	variables["user_id"] = userID
//
//	return c.transport.Raw(ctx, `
//		mutation UnassignRole($role_id: ID!, $user_id: String!) {
//			unassignRole(input: { role_id: $role_id, user_id: $user_id })
//		}
//	`, variables, nil)
//}
//
//// ListRoles returns all of the roles in the database
//func (c *Client) ListRoles(ctx context.Context) ([]*primitives.Role, error) {
//	var resp struct {
//		Roles []*primitives.Role `json:"roles"`
//	}
//
//	err := c.transport.Raw(ctx, `
//		query Roles {
//			roles {
//				id
//				label
//				system
//			}
//		}
//	`, nil, &resp)
//
//	if err != nil {
//		return nil, err
//	}
//
//	return resp.Roles, nil
//}

// GetUsers returns all users for the given emails
func (c *Client) GetUsers(ctx context.Context, emails []primitives.Email) ([]*models.User, error) {
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
	Secret primitives.Password `json:"secret"`
	Token  *primitives.Token   `json:"token"`
}

type CreateTokenResponse struct {
	Response *CreateTokenMutation `json:"createToken"`
}

// CreateToken creates a new API token for the provided user. You can pass nil and it will return a token for you
func (c *Client) CreateToken(ctx context.Context, user *models.User) (*auth.APIToken, *primitives.Token, error) {
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
	IDs []database.ID `json:"tokens"`
}

// ListTokens lists all of the auth tokens for the provided user
func (c *Client) ListTokens(ctx context.Context, user *models.User) ([]database.ID, error) {
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
func (c *Client) RemoveToken(ctx context.Context, tokenID database.ID) error {
	variables := make(map[string]interface{})
	variables["id"] = tokenID

	return c.transport.Raw(ctx, `
		mutation RemoveToken($id: ID!) {
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

type GetProjectResponse struct {
	*models.Project
	Contributors []GQLContributor `json:"contributors"`
}

func (c *Client) GetProject(ctx context.Context, id string, label *models.Label) (*GetProjectResponse, error) {
	var resp struct {
		Project GetProjectResponse `json:"project"`
	}

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

				contributors {
					id
					user {
						id
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

	return &resp.Project, nil
}

type UpdateProjectSpecResponseBody struct {
	*models.Project
	ProjectSpec *models.ProjectSpec `json:"current_spec"`
}

type UpdateProjectSpecResponse struct {
	UpdateProjectSpecResponseBody `json:"updateProjectSpec"`
}

func (c *Client) UpdateProjectSpec(ctx context.Context, projectLabel models.Label, spec *models.ProjectSpecFile) (*models.Project, *models.ProjectSpec, error) {
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
					policy
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

func (c *Client) AddContributor(ctx context.Context, project models.Project, user models.User, role models.Label) (*models.Contributor, error) {
	variables := map[string]interface{}{
		"project_label": project.Label,
		"email":         user.Email,
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

func (c *Client) AttemptRecovery(ctx context.Context, ID string, secret primitives.Password, newPassword primitives.Password) error {
	variables := map[string]interface{}{
		"new_password": newPassword.String(),
		"secret":       secret.String(),
		"id":           ID,
	}

	query := `mutation attemptRecovery($new_password: Password!, $secret: Password!, $id: ID!) {
		attemptRecovery(input: {
			new_password: $new_password,
			secret: $secret,
			id: $id,
		})
	}`

	return c.transport.Raw(ctx, query, variables, nil)
}

type ListRecoveriesResponse struct {
	Recoveries []*primitives.Recovery `json:"recoveries"`
}

func (c *Client) Recoveries(ctx context.Context) ([]*primitives.Recovery, error) {
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

func (c *Client) DeleteRecoveries(ctx context.Context, ids []database.ID) error {
	variables := map[string]interface{}{
		"ids": ids,
	}

	query := `
		mutation DeleteRecoveries($ids: [ID!]!) {
			deleteRecoveries(input: { ids: $ids })
		}
	`

	return c.transport.Raw(ctx, query, variables, nil)
}
