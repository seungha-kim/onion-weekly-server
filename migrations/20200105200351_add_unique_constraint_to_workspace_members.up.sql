alter table workspace_members
    add constraint workspace_members_user_id_workspace_id_unique unique (workspace_id, user_id);