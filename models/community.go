package models

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"time"

	"glorp/config"
)

type Community struct {
	ID            int                   `json:"id"`
	Name          string                `json:"name"`
	DisplayName   string                `json:"display_name"`
	Description   string                `json:"description"`
	CreatorID     int                   `json:"creator_id"`
	Creator       *User                 `json:"creator,omitempty"`
	Visibility    string                `json:"visibility"`    // public, private, restricted
	JoinApproval  string                `json:"join_approval"` // open, approval_required, invite_only
	MemberCount   int                   `json:"member_count"`
	UserRole      string                `json:"user_role,omitempty"`      // Current user's role in community
	UserStatus    string                `json:"user_status,omitempty"`    // Current user's membership status
	JoinRequested bool                  `json:"join_requested,omitempty"` // Whether user has pending request
	CreatedAt     time.Time             `json:"created_at"`
	UpdatedAt     time.Time             `json:"updated_at"`
	Rules         []CommunityRule       `json:"rules,omitempty"`
	Moderators    []CommunityMembership `json:"moderators,omitempty"`
}

type CommunityMembership struct {
	ID          int       `json:"id"`
	CommunityID int       `json:"community_id"`
	UserID      int       `json:"user_id"`
	User        *User     `json:"user,omitempty"`
	Role        string    `json:"role"`   // member, moderator, admin, creator
	Status      string    `json:"status"` // active, pending, banned
	JoinedAt    time.Time `json:"joined_at"`
}

type CommunityJoinRequest struct {
	ID          int        `json:"id"`
	CommunityID int        `json:"community_id"`
	Community   *Community `json:"community,omitempty"`
	UserID      int        `json:"user_id"`
	User        *User      `json:"user,omitempty"`
	Message     string     `json:"message"`
	Status      string     `json:"status"` // pending, approved, rejected
	RequestedAt time.Time  `json:"requested_at"`
	ProcessedAt *time.Time `json:"processed_at,omitempty"`
	ProcessedBy *int       `json:"processed_by,omitempty"`
	Processor   *User      `json:"processor,omitempty"`
}

type CommunityRule struct {
	ID          int       `json:"id"`
	CommunityID int       `json:"community_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	RuleOrder   int       `json:"rule_order"`
	CreatedAt   time.Time `json:"created_at"`
}

type CommunityFilters struct {
	Search     string
	Visibility string
	SortBy     string
	Filter     string
	UserID     int
	Page       int
	Limit      int
}

// Community CRUD operations
func CreateCommunity(name, displayName, description string, creatorID int, visibility, joinApproval string) (*Community, error) {
	// Validate community name
	if err := validateCommunityName(name); err != nil {
		return nil, err
	}

	// Check if community name already exists
	var count int
	err := config.DB.QueryRow("SELECT COUNT(*) FROM communities WHERE name = ?", name).Scan(&count)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, fmt.Errorf("community name already exists")
	}

	tx, err := config.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Create community
	result, err := tx.Exec(`
		INSERT INTO communities (name, display_name, description, creator_id, visibility, join_approval) 
		VALUES (?, ?, ?, ?, ?, ?)`,
		name, displayName, description, creatorID, visibility, joinApproval)
	if err != nil {
		return nil, err
	}

	communityID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Add creator as community admin
	_, err = tx.Exec(`
		INSERT INTO community_memberships (community_id, user_id, role, status) 
		VALUES (?, ?, 'creator', 'active')`,
		communityID, creatorID)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return GetCommunityByID(int(communityID), creatorID)
}

func GetCommunityByID(id, userID int) (*Community, error) {
	community := &Community{}
	query := `
		SELECT c.id, c.name, c.display_name, c.description, c.creator_id, c.visibility, 
		       c.join_approval, c.member_count, c.created_at, c.updated_at, u.username
		FROM communities c
		LEFT JOIN users u ON c.creator_id = u.id
		WHERE c.id = ?`

	var creatorUsername string
	err := config.DB.QueryRow(query, id).Scan(
		&community.ID, &community.Name, &community.DisplayName, &community.Description,
		&community.CreatorID, &community.Visibility, &community.JoinApproval,
		&community.MemberCount, &community.CreatedAt, &community.UpdatedAt, &creatorUsername)

	if err != nil {
		return nil, err
	}

	community.Creator = &User{
		ID:       community.CreatorID,
		Username: creatorUsername,
	}

	// Get user's role and status in community if userID provided
	if userID > 0 {
		var role, status sql.NullString
		memberQuery := `SELECT role, status FROM community_memberships WHERE community_id = ? AND user_id = ?`
		config.DB.QueryRow(memberQuery, community.ID, userID).Scan(&role, &status)

		if role.Valid {
			community.UserRole = role.String
			community.UserStatus = status.String
		}

		// Check if user has pending join request
		var requestCount int
		requestQuery := `SELECT COUNT(*) FROM community_join_requests WHERE community_id = ? AND user_id = ? AND status = 'pending'`
		config.DB.QueryRow(requestQuery, community.ID, userID).Scan(&requestCount)
		community.JoinRequested = requestCount > 0
	}

	return community, nil
}

func GetCommunityByName(name string, userID int) (*Community, error) {
	community := &Community{}
	query := `
		SELECT c.id, c.name, c.display_name, c.description, c.creator_id, c.visibility, 
		       c.join_approval, c.member_count, c.created_at, c.updated_at, u.username
		FROM communities c
		LEFT JOIN users u ON c.creator_id = u.id
		WHERE c.name = ?`

	var creatorUsername string
	err := config.DB.QueryRow(query, name).Scan(
		&community.ID, &community.Name, &community.DisplayName, &community.Description,
		&community.CreatorID, &community.Visibility, &community.JoinApproval,
		&community.MemberCount, &community.CreatedAt, &community.UpdatedAt, &creatorUsername)

	if err != nil {
		return nil, err
	}

	community.Creator = &User{
		ID:       community.CreatorID,
		Username: creatorUsername,
	}

	// Get user's role and status in community if userID provided
	if userID > 0 {
		var role, status sql.NullString
		memberQuery := `SELECT role, status FROM community_memberships WHERE community_id = ? AND user_id = ?`
		config.DB.QueryRow(memberQuery, community.ID, userID).Scan(&role, &status)

		if role.Valid {
			community.UserRole = role.String
			community.UserStatus = status.String
		}

		// Check if user has pending join request
		var requestCount int
		requestQuery := `SELECT COUNT(*) FROM community_join_requests WHERE community_id = ? AND user_id = ? AND status = 'pending'`
		config.DB.QueryRow(requestQuery, community.ID, userID).Scan(&requestCount)
		community.JoinRequested = requestCount > 0
	}

	return community, nil
}

func GetCommunities(filters CommunityFilters) ([]Community, int, error) {
	var communities []Community
	var whereConditions []string
	var args []interface{}

	baseQuery := `
		FROM communities c 
		LEFT JOIN users u ON c.creator_id = u.id
	`

	// Build WHERE conditions
	if filters.Search != "" {
		whereConditions = append(whereConditions, "(c.name LIKE ? OR c.display_name LIKE ? OR c.description LIKE ?)")
		searchTerm := "%" + filters.Search + "%"
		args = append(args, searchTerm, searchTerm, searchTerm)
	}

	if filters.Visibility != "" {
		whereConditions = append(whereConditions, "c.visibility = ?")
		args = append(args, filters.Visibility)
	}

	// Handle filter type
	switch filters.Filter {
	case "joined":
		if filters.UserID > 0 {
			whereConditions = append(whereConditions, `c.id IN (
				SELECT community_id FROM community_memberships 
				WHERE user_id = ? AND status = 'active'
			)`)
			args = append(args, filters.UserID)
		}
	case "popular":
		// No additional WHERE conditions needed, will be handled by ORDER BY
	case "recent":
		// No additional WHERE conditions needed, will be handled by ORDER BY
	}

	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Count total
	countQuery := "SELECT COUNT(*) " + baseQuery + whereClause
	var total int
	err := config.DB.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Main query
	selectQuery := `
		SELECT c.id, c.name, c.display_name, c.description, c.creator_id, c.visibility, 
		       c.join_approval, c.member_count, c.created_at, c.updated_at, u.username
	` + baseQuery + whereClause

	// Add sorting based on filter type
	switch filters.SortBy {
	case "popular":
		// Hot algorithm: member_count / (age in hours + 2)^1.8
		selectQuery += " ORDER BY (c.member_count + 1) / POW((julianday('now') - julianday(c.created_at)) * 24 + 2, 1.8) DESC"
	case "top":
		selectQuery += " ORDER BY c.member_count DESC, c.created_at DESC"
	case "new":
		selectQuery += " ORDER BY c.created_at DESC"
	default: // date or popularity
		selectQuery += " ORDER BY c.created_at DESC"
	}

	// Add pagination
	if filters.Limit > 0 {
		selectQuery += " LIMIT ?"
		args = append(args, filters.Limit)

		if filters.Page > 0 {
			offset := (filters.Page - 1) * filters.Limit
			selectQuery += " OFFSET ?"
			args = append(args, offset)
		}
	}

	rows, err := config.DB.Query(selectQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var community Community
		var creatorUsername string

		err := rows.Scan(
			&community.ID, &community.Name, &community.DisplayName, &community.Description,
			&community.CreatorID, &community.Visibility, &community.JoinApproval,
			&community.MemberCount, &community.CreatedAt, &community.UpdatedAt, &creatorUsername,
		)
		if err != nil {
			continue
		}

		community.Creator = &User{
			ID:       community.CreatorID,
			Username: creatorUsername,
		}

		communities = append(communities, community)
	}

	return communities, total, nil
}

// Membership operations
func JoinCommunity(communityID, userID int, message string) error {
	community, err := GetCommunityByID(communityID, userID)
	if err != nil {
		return err
	}

	// Check if user is already a member
	if community.UserRole != "" {
		return fmt.Errorf("already a member of this community")
	}

	// Check if user has pending request
	if community.JoinRequested {
		return fmt.Errorf("join request already pending")
	}

	switch community.JoinApproval {
	case "open":
		// Directly add as member
		return addCommunityMember(communityID, userID, "member")
	case "approval_required":
		// Create join request
		return createJoinRequest(communityID, userID, message)
	case "invite_only":
		return fmt.Errorf("this community is invite-only")
	default:
		return fmt.Errorf("invalid join approval setting")
	}
}

func LeaveCommunity(communityID, userID int) error {
	// Check if user is creator
	var creatorID int
	err := config.DB.QueryRow("SELECT creator_id FROM communities WHERE id = ?", communityID).Scan(&creatorID)
	if err != nil {
		return err
	}

	if creatorID == userID {
		return fmt.Errorf("creator cannot leave community")
	}

	tx, err := config.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Remove membership
	_, err = tx.Exec("DELETE FROM community_memberships WHERE community_id = ? AND user_id = ?", communityID, userID)
	if err != nil {
		return err
	}

	// Update member count
	_, err = tx.Exec("UPDATE communities SET member_count = member_count - 1 WHERE id = ?", communityID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func ProcessJoinRequest(requestID, processorID int, approved bool) error {
	// Get request details
	var communityID, userID int
	var status string
	err := config.DB.QueryRow("SELECT community_id, user_id, status FROM community_join_requests WHERE id = ?", requestID).Scan(&communityID, &userID, &status)
	if err != nil {
		return err
	}

	if status != "pending" {
		return fmt.Errorf("request already processed")
	}

	// Check if processor has permission
	if !CanManageCommunity(communityID, processorID) {
		return fmt.Errorf("insufficient permissions")
	}

	tx, err := config.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	newStatus := "rejected"
	if approved {
		newStatus = "approved"
		// Add user as member
		_, err = tx.Exec(`
			INSERT INTO community_memberships (community_id, user_id, role, status) 
			VALUES (?, ?, 'member', 'active')`,
			communityID, userID)
		if err != nil {
			return err
		}

		// Update member count
		_, err = tx.Exec("UPDATE communities SET member_count = member_count + 1 WHERE id = ?", communityID)
		if err != nil {
			return err
		}
	}

	// Update request
	_, err = tx.Exec(`
		UPDATE community_join_requests 
		SET status = ?, processed_at = CURRENT_TIMESTAMP, processed_by = ? 
		WHERE id = ?`,
		newStatus, processorID, requestID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// Helper functions
func validateCommunityName(name string) error {
	if len(name) < 3 {
		return fmt.Errorf("community name must be at least 3 characters")
	}
	if len(name) > 50 {
		return fmt.Errorf("community name must be less than 50 characters")
	}

	// Check for valid characters (alphanumeric and underscores only)
	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(name) {
		return fmt.Errorf("community name can only contain letters, numbers, and underscores")
	}

	return nil
}

func addCommunityMember(communityID, userID int, role string) error {
	tx, err := config.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Add membership
	_, err = tx.Exec(`
		INSERT INTO community_memberships (community_id, user_id, role, status) 
		VALUES (?, ?, ?, 'active')`,
		communityID, userID, role)
	if err != nil {
		return err
	}

	// Update member count
	_, err = tx.Exec("UPDATE communities SET member_count = member_count + 1 WHERE id = ?", communityID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func createJoinRequest(communityID, userID int, message string) error {
	_, err := config.DB.Exec(`
		INSERT INTO community_join_requests (community_id, user_id, message) 
		VALUES (?, ?, ?)`,
		communityID, userID, message)
	return err
}

func CanManageCommunity(communityID, userID int) bool {
	var role string
	err := config.DB.QueryRow(`
		SELECT role FROM community_memberships 
		WHERE community_id = ? AND user_id = ? AND status = 'active'`,
		communityID, userID).Scan(&role)

	if err != nil {
		return false
	}

	return role == "creator" || role == "admin" || role == "moderator"
}

func GetCommunityModerators(communityID int) ([]CommunityMembership, error) {
	var moderators []CommunityMembership
	query := `
		SELECT cm.id, cm.community_id, cm.user_id, cm.role, cm.status, cm.joined_at,
		       u.username, u.display_name
		FROM community_memberships cm
		LEFT JOIN users u ON cm.user_id = u.id
		WHERE cm.community_id = ? AND cm.role IN ('creator', 'admin', 'moderator') AND cm.status = 'active'
		ORDER BY 
			CASE cm.role 
				WHEN 'creator' THEN 1 
				WHEN 'admin' THEN 2 
				WHEN 'moderator' THEN 3 
			END, cm.joined_at ASC`

	rows, err := config.DB.Query(query, communityID)
	if err != nil {
		return moderators, err
	}
	defer rows.Close()

	for rows.Next() {
		var membership CommunityMembership
		var username, displayName string

		err := rows.Scan(
			&membership.ID, &membership.CommunityID, &membership.UserID,
			&membership.Role, &membership.Status, &membership.JoinedAt,
			&username, &displayName,
		)
		if err != nil {
			continue
		}

		membership.User = &User{
			ID:          membership.UserID,
			Username:    username,
			DisplayName: displayName,
		}

		moderators = append(moderators, membership)
	}

	return moderators, nil
}

func GetPendingJoinRequests(communityID int) ([]CommunityJoinRequest, error) {
	var requests []CommunityJoinRequest
	query := `
		SELECT r.id, r.community_id, r.user_id, r.message, r.status, r.requested_at,
		       u.username, u.display_name
		FROM community_join_requests r
		LEFT JOIN users u ON r.user_id = u.id
		WHERE r.community_id = ? AND r.status = 'pending'
		ORDER BY r.requested_at ASC`

	rows, err := config.DB.Query(query, communityID)
	if err != nil {
		return requests, err
	}
	defer rows.Close()

	for rows.Next() {
		var request CommunityJoinRequest
		var username, displayName string

		err := rows.Scan(
			&request.ID, &request.CommunityID, &request.UserID,
			&request.Message, &request.Status, &request.RequestedAt,
			&username, &displayName,
		)
		if err != nil {
			continue
		}

		request.User = &User{
			ID:          request.UserID,
			Username:    username,
			DisplayName: displayName,
		}

		requests = append(requests, request)
	}

	return requests, nil
}

// GetCommunityCount returns the total number of communities
func GetCommunityCount() (int, error) {
	var count int
	err := config.DB.QueryRow("SELECT COUNT(*) FROM communities").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// UpdateCommunity updates a community's name and description
func UpdateCommunity(communityID int, name, description string) error {
	query := "UPDATE communities SET name = ?, description = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?"
	_, err := config.DB.Exec(query, name, description, communityID)
	return err
}

// DeleteCommunity deletes a community and all its associated data
func DeleteCommunity(communityID int) error {
	// Start a transaction
	tx, err := config.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete community memberships
	_, err = tx.Exec("DELETE FROM community_memberships WHERE community_id = ?", communityID)
	if err != nil {
		return err
	}

	// Delete community join requests
	_, err = tx.Exec("DELETE FROM community_join_requests WHERE community_id = ?", communityID)
	if err != nil {
		return err
	}

	// Delete community rules
	_, err = tx.Exec("DELETE FROM community_rules WHERE community_id = ?", communityID)
	if err != nil {
		return err
	}

	// Update threads to remove community association
	_, err = tx.Exec("UPDATE threads SET community_id = NULL WHERE community_id = ?", communityID)
	if err != nil {
		return err
	}

	// Finally, delete the community
	_, err = tx.Exec("DELETE FROM communities WHERE id = ?", communityID)
	if err != nil {
		return err
	}

	// Commit the transaction
	return tx.Commit()
}

// TransferCommunityOwnership transfers ownership of a community to a new user
func TransferCommunityOwnership(communityID, newOwnerID int) error {
	// Start a transaction
	tx, err := config.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update community creator
	_, err = tx.Exec("UPDATE communities SET creator_id = ? WHERE id = ?", newOwnerID, communityID)
	if err != nil {
		return err
	}

	// Update community membership for new owner
	_, err = tx.Exec(`
		INSERT INTO community_memberships (community_id, user_id, role, status, joined_at)
		VALUES (?, ?, 'creator', 'active', CURRENT_TIMESTAMP)
		ON CONFLICT(community_id, user_id) DO UPDATE SET role = 'creator', status = 'active'
	`, communityID, newOwnerID)
	if err != nil {
		return err
	}

	// Update old creator's role to member
	_, err = tx.Exec(`
		UPDATE community_memberships 
		SET role = 'member' 
		WHERE community_id = ? AND role = 'creator'
	`, communityID)
	if err != nil {
		return err
	}

	// Commit the transaction
	return tx.Commit()
}
