# A Users's Guide to Thunderdome

Thunderdome is a fun way to facilitate agile scrum practices including Story pointing (games), Sprint Retrospectives,
Story mapping, Async Daily Standup (team standup) and more.

As a new user in this realm, let this be your guide. First, we need to know who you are.

## Register (optional)

Create a new account, or join as guest.

![image](https://github.com/scottux/thunderdome-planning-poker/assets/135884/c6f7a12c-4d81-4622-bae4-9356bdf16e62)

Having an account lets you save your games and more.

![image](https://github.com/scottux/thunderdome-planning-poker/assets/135884/395d06ca-6485-4293-8482-31a685e3db56)

- Name  
  This will be visible to others.
- Email (for account)
- Password (for account)  
  Use a strong password. Type it again to confirm.

You will receive an email to confirm your new account.

## Login

Use the email/password you created when registering.

![image](https://github.com/scottux/thunderdome-planning-poker/assets/135884/7a9ff086-fd30-4496-a1ee-a083370b9a5d)

- Email
- Password

_OIDC Providers coming soon._

### Password Retrieval

Forgot your password? Thunderdome can send a password reset link to your email.

![image](https://github.com/scottux/thunderdome-planning-poker/assets/135884/a42e8620-dbe8-4ec4-a784-63d19fcc1327)

## Profile

User, it is all about you! Control your Thunderdome experience.

### Details

![image](https://github.com/scottux/thunderdome-planning-poker/assets/135884/4c4a4512-538e-4150-b0cf-7864046fb1ff)

- Name  
  This will be visible to others.
- Email  
  Update your account email, or enter one if you are a guest.
- Country (optional)
- Locale, default: English  
  8 locales to choose from.
- Company (optional)
- Job Title (optional)
- Theme, default: auto  
  The default lets the operating system and browser define dark or light, if supported. If you prefer a darker or
  lighter interface, you may choose it here.
- Option: Enable Game Notification, default: true
- Avatar, default: robohash  
  Several others to choose from, pick your flavor; mp, identicon, monsterid, wavatar, retro.

### API Access

Create an API key to integrate Thunderdome with your tools.

See API Documentation here [Thunderdome API Docs](https://thunderdome.dev/swagger/index.html)

![image](https://github.com/scottux/thunderdome-planning-poker/assets/135884/431ef5a7-04fc-408c-90a7-f3dbb8690bf8)

![image](https://github.com/scottux/thunderdome-planning-poker/assets/135884/8e81b994-1cca-4f84-8bbb-4957527573fb)

### Jira Integration (premium only)

Integrate directly with your team's backlog to import your stories to point.

_Other integrations coming soon._

### Delete Account

This is the Thunderdome, but you are free to leave. We will erase all data about you.

**This is permanent:** All poker sessions, retros, story maps, and orgs/teams directly owned by your account will also
be deleted.

![image](https://github.com/scottux/thunderdome-planning-poker/assets/135884/ae903d0a-d58c-4319-bf62-455e4ebf040f)

## Game

In Thunderdome, an agile poker planning session is known as a Game.

You can create a game to determine the size of a story, or join one in progress.

### Create a Game

![image](https://github.com/scottux/thunderdome-planning-poker/assets/135884/da9b6c0b-e0ba-48d6-b999-4dca19772f5f)

- Name
- Team (optional)
- Point Range Allowed, default: [ 1, 2, 3, 5, 8, 13, ? ]
- Stories  
  Upload an XML or CSV for stories, or add manually. See note in Stories.
- Point Average Rounding, default: Ceil  
  Other options; Round, Floor.
- Option: Auto Finish Voting, default: true
- Option: Hide Voter Identity, default: false
- Passcode (optional)
- Leader Code (optional)

### Game

The planning session is real-time, each user chooses the size for the story and votes are shown when everyone has
finished.

If the team agrees, the game is over. If not, then it has just begun!

![image](https://user-images.githubusercontent.com/846933/95778842-eb76ef00-0c96-11eb-99d8-af5d098c12ee.png)

### Stories

This can be a list of Stories, Bugs, Tasks, etc. and serves as a queue for team voting.

#### Import stories from Jira Cloud

Premium feature.

#### Import stories from Jira XML

Upload.

#### Import stories from a CSV file

The CSV file must include all the following fields with no header row:

_Type,Title,ReferenceId,Link,Description,AcceptanceCriteria_

#### Create a Story

- Type, default: Story  
  Other Types: Bug, Spike, Epic, Task, Subtask
- Name
- Reference ID
- Link
- Priority  
  Priorities; Blocker, Highest, High, Medium, Low, Lowest
- Description  
  A full text editor is supplied to provide a detailed description.
- Acceptance Criteria  
  A full text editor is supplied, feel free to use Gherkin statements.

### Users

See who is voting or become a spectator.

### Invite

Send a link for others to join the game.

## Retrospectives

Facilitates an agile sprint retrospective.

Retrospectives happen in phases. The first phase is the Prime Directive. You may edit or delete the retro at any time.

1. Prime Directive
2. Brainstorm  
   Add comments. What went well? What needs improvement? I want to ask...
3. Group  
   Organize comments into topics. Drag and drop to sort.
4. Vote  
   Vote for groups to discuss first.
5. Action Items  
   Add Action Items, the grouping and voting phases become locked.
6. Done  
   Export the Retro

![image](https://user-images.githubusercontent.com/846933/173260209-3ef3299f-f1b2-41e8-802f-17d40649c66d.png)

### Create a Retro

- Name
- Team (optional)
- Join Code (optional)
- Fac. Code (optional)
- Max Group Votes per User, default: 3
- Brainstorm Phase Feedback Visibility, default: Feedback Visible  
  Other options include concealed and hidden. Determines if team members can see each other's suggestions.

## Storyboards

Stories are units of work that need to be sized.

### Goals

A goal is a way to group stories.

- Name (optional)

#### Columns

A column is customizable and serves as a way to track stories throughout the goal.

- Title Text (optional)

#### Add Story

A story is a unit of work. It can be in an open or closed state.

- Name
- Link
- Points
- Color
- Content
- Discussion

### Personas

_Coming Soon_

### Add Persona

- Name
- Role
- Description

### Color Legend

A palette is provided so that you can choose to apply meaningful colors to story cards.

#### Edit Legend

This is where you can define what each color means.

![image](https://user-images.githubusercontent.com/846933/173260211-304a973d-4ede-494f-bb7d-b7e5c86a4e6e.png)

### Create a Storyboard

- Name
- Team (optional)
- Passcode (optional)
- Facilitator Code (optional)

## Teams, Organizations, and Departments

![image](https://github.com/StevenWeathers/thunderdome-planning-poker/assets/135884/9a3cc4ea-3a6b-4b00-9664-1941209d7568)

Teams can be simple, or they can be within Organizations and Departments.

### Organizations

![image](https://github.com/StevenWeathers/thunderdome-planning-poker/assets/135884/451a6ba4-c712-4e41-8928-bce48c942e20)

#### Create Organization

- Name

### Departments

![image](https://github.com/StevenWeathers/thunderdome-planning-poker/assets/135884/8fdb047b-8abf-4bfe-b8c1-5904ece4934b)

#### Create Department

- Name

### Teams

#### Create Team

- Name

#### Add User

- User Email
- Role  
  Admin or Member

#### Checkins

Asynchronous daily standup tool to aid in speeding up standup or making standups completely async depending on team
practices.

![image](https://github.com/StevenWeathers/thunderdome-planning-poker/assets/135884/eab485f4-7254-46bf-b2e8-11e468ba47f0)

##### Check In

Provide your daily standup report. What did you do yesterday? What are you doing today? Any blockers? Anything to
discuss?

![image](https://github.com/StevenWeathers/thunderdome-planning-poker/assets/135884/e4411409-421e-4c54-bd4e-d73ab006153a)

Choose your timezone.

![image](https://github.com/StevenWeathers/thunderdome-planning-poker/assets/135884/00af5a6c-f19d-4de9-916a-a2cec975f60d)

##### Create Games, Retros, and Storyboards

Creating these sessions within a Team, Organization, or Department will pre-populate those fields.

See each individual section for further details about creating games, retros, and storyboards.

## Languages

🌍 Thunderdome has made every effort to be an international tool, and help developers of all nationalities unite
together.

## Contributions

Thunderdome is released as open source software, the code is hosted on Github and licensed Apache 2.0.

_v3.6.3_
