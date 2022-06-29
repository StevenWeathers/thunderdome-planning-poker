DROP TABLE retro_action_comment;
DROP TABLE retro_action_assignee;

DROP PROCEDURE retro_action_comment_add(UUID, UUID, UUID, TEXT);
DROP PROCEDURE retro_action_comment_edit(UUID, UUID, UUID, TEXT);
DROP PROCEDURE retro_action_comment_delete(UUID, UUID, UUID);
DROP PROCEDURE retro_action_assignee_add(UUID, UUID, UUID);
DROP PROCEDURE retro_action_assignee_delete(UUID, UUID, UUID);