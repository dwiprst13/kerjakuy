package rbac

type Role string
type Permission string

const (
	RoleOwner  Role = "owner"
	RoleAdmin  Role = "admin"
	RoleMember Role = "member"
)

const (
	// Workspace permissions
	PermissionUpdateWorkspace Permission = "workspace:update"
	PermissionDeleteWorkspace Permission = "workspace:delete"
	PermissionInviteMember    Permission = "workspace:invite_member"
	PermissionRemoveMember    Permission = "workspace:remove_member"
	PermissionUpdateMember    Permission = "workspace:update_member"

	// Project permissions
	PermissionCreateProject Permission = "project:create"
	PermissionUpdateProject Permission = "project:update"
	PermissionDeleteProject Permission = "project:delete"

	// Board/Column/Task permissions 
	PermissionCreateBoard  Permission = "board:create"
	PermissionUpdateBoard  Permission = "board:update"
	PermissionDeleteBoard  Permission = "board:delete"
	PermissionCreateTask   Permission = "task:create"
	PermissionUpdateTask   Permission = "task:update"
	PermissionDeleteTask   Permission = "task:delete"
)

var Policy = map[Role][]Permission{
	RoleOwner: {
		PermissionUpdateWorkspace,
		PermissionDeleteWorkspace,
		PermissionInviteMember,
		PermissionRemoveMember,
		PermissionUpdateMember,
		PermissionCreateProject,
		PermissionUpdateProject,
		PermissionDeleteProject,
		PermissionCreateBoard,
		PermissionUpdateBoard,
		PermissionDeleteBoard,
		PermissionCreateTask,
		PermissionUpdateTask,
		PermissionDeleteTask,
	},
	RoleAdmin: {
		PermissionUpdateWorkspace,
		PermissionInviteMember,
		PermissionRemoveMember,
		PermissionUpdateMember,
		PermissionCreateProject,
		PermissionUpdateProject,
		PermissionDeleteProject,
		PermissionCreateBoard,
		PermissionUpdateBoard,
		PermissionDeleteBoard,
		PermissionCreateTask,
		PermissionUpdateTask,
		PermissionDeleteTask,
	},
	RoleMember: {
		PermissionCreateProject, 
		PermissionCreateProject,
		PermissionUpdateProject,
		PermissionCreateBoard,
		PermissionUpdateBoard,
		PermissionCreateTask,
		PermissionUpdateTask,
		PermissionDeleteTask,
	},
}

func HasPermission(role Role, perm Permission) bool {
	perms, ok := Policy[role]
	if !ok {
		return false
	}
	for _, p := range perms {
		if p == perm {
			return true
		}
	}
	return false
}
