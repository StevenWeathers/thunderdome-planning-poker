# Socket Events

## The Server

### Emits

- **'init'** The connection is established, current state is passed over.
- **'user_activity'** Someone has left or joined, but not left inner joined.
- **'vote_activity'** Someone voted.
- **'chat'** Passed through from a client.
- **'logged_in'** The user has supplied the proper credentials.
- **'session_list'** The overlord has requested a list of their sessions.
- **'task_list'** The overlord has requested a list of their tasks for a session.
- **'task_init'** The overlord has chosen a task, all clients are now aware.
- **'vote_init'** The overlord has decided to begin voting. Clients can choose numbers.
- **'vote_end'** The voting session has stopped, await the overlord's next move.
- **'task_end'** This task is sized and notes are captured, we are moving on.

### Reacts to

- **'chat'** Passes it through to all clients.
- **'vote'** Accepts the user and their vote for the task, emits vote_activity in response.
- **'get_tasks'** If an overlord, fetches the tasks. Emits a task_list back to the overlord.
- **'task_start'** Accepts the task and emits a task_init to all clients.
- **'start_sizing'** Accepts any updates to the task notes, emits vote_init to all clients.
- **'stop_sizing'** Accepts any updates to the task notes, emits vote_end to all clients.
- **'task_stop'** Saves the task and emits a task_end to the sender.

## The Client

### Emits

- **'chat'** A user is chatting, contains a message.
- **'vote'** A user is voting, contains a size.
- **'get_tasks'** The overlord is requesting sizes, contains the planning session id.
- **'task_start'** The overlord is choosing a task to size, sends the task id.
- **'start_sizing'** The overlord has pressed the "Begin Sizing" button.
- **'stop_sizing'** The overlord has pressed the "Stop Sizing" button.
- **'task_stop'** The overlord has pressed the "Finished with this task" button.

### Reacts to

- **'init'** Determine if we are in a session already, or if there is a task already chosen.
- **'user_activity'** Refresh the users list.
- **'vote_activity'** Refresh the voters list.
- **'chat'** Update the chat window.
- **'logged_in'** isOverlord?
- **'session_list'** Display the list of sessions to the overlord.
- **'task_list'**