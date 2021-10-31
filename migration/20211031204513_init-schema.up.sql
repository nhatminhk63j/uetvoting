create table if not exists users
(
  id         int          not null primary key auto_increment,
  name       varchar(255) not null,
  email      varchar(255) not null,
  picture    varchar(255) not null,
  role       int          not null,
  created_at timestamp default current_timestamp,
  updated_at timestamp default current_timestamp on update current_timestamp
) engine = InnoDB;

create table if not exists events
(
  id           int          not null primary key auto_increment,
  name         varchar(255) not null,
  participants text         not null,
  is_approved  bool      default false,
  is_opened    bool      default false,
  created_by   int          not null,
  updated_by   int          not null,
  created_at   timestamp default current_timestamp,
  updated_at   timestamp default current_timestamp on update current_timestamp,
  constraint `fk_events_users_1`
    foreign key (created_by) references users (id)
      on update no action
      on delete no action,
  constraint `fk_events_users_2`
    foreign key (updated_by) references users (id)
      on update no action
      on delete no action
) engine = InnoDB;

create table if not exists sections
(
  id          int          not null primary key auto_increment,
  name        varchar(255) not null,
  description text,
  position    int          not null,
  event_id    int          not null,
  created_at  timestamp default current_timestamp,
  updated_at  timestamp default current_timestamp on update current_timestamp,
  constraint `fk_sections_events`
    foreign key (event_id) references events (id)
      on update no action
      on delete no action
) engine = InnoDB;

create table if not exists questions
(
  id             int               not null primary key auto_increment,
  content        text              not null,
  position       int               not null,
  type           enum ('CHECKBOX') not null,
  is_required    bool      default false,
  limited_choice int               not null,
  section_id     int               not null,
  created_at     timestamp default current_timestamp,
  updated_at     timestamp default current_timestamp on update current_timestamp,
  constraint `fk_questions_sections`
    foreign key (section_id) references sections (id)
      on delete no action
      on update no action
) engine = InnoDB;

create table if not exists options
(
  id          int  not null primary key auto_increment,
  content     text not null,
  question_id int  not null,
  constraint `fk_options_questions`
    foreign key (question_id) references questions (id)
      on update no action
      on delete no action
) engine = InnoDB;

create table if not exists votes
(
  id           varchar(255) not null primary key default (uuid()),
  is_completed bool                              default false,
  event_id     int          not null,
  constraint `fk_votes_events`
    foreign key (event_id) references events (id)
      on update no action
      on delete no action
) engine = InnoDB;

create table if not exists answers
(
  question_id int          not null,
  vote_id     varchar(255) not null,
  content     text,
  created_at  timestamp default current_timestamp,
  updated_at  timestamp default current_timestamp on update current_timestamp,
  primary key (question_id, vote_id),
  constraint `fk_answers_questions`
    foreign key (question_id) references questions (id)
      on update no action
      on delete no action,
  constraint `fk_answers_votes`
    foreign key (vote_id) references votes (id)
      on update no action
      on delete no action
) engine = InnoDB;
