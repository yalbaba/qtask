// +build !oracle

package sql

const SQLGetSEQ = `insert into tsk_system_seq (name) values (@name)`

const SQLCreateTaskID = `insert into tsk_system_task
(task_id,
 name,
 next_execute_time,
 max_execute_time,
 next_interval,
 delete_interval,
 status,
 queue_name,
 msg_content)
values
(@task_id,
 @name, 
 date_add(now(),interval #first_timeout second),   
 date_add(now(),interval #max_timeout second),
 @next_interval,
 @delete_interval,
 20,
 @queue_name,
 @content)`

const SQLProcessingTask = `
update tsk_system_task t set 
t.next_execute_time=date_add(now(),interval t.next_interval second),
t.status=30,
t.count=t.count + 1,
t.last_execute_time=now()
where t.task_id=@task_id 
and t.status in(20,30)
and t.count <= 5`

const SQLFinishTask = `
update tsk_system_task t
set t.next_execute_time = STR_TO_DATE('2099-12-31', '%Y-%m-%d'),
    t.status            = 0,
    t.delete_time       = date_add(now(),interval t.delete_interval second)
where t.task_id = @task_id
and t.status in (20, 30)`

const SQLUpdateTask = `
update tsk_system_task t set 
t.batch_id=@batch_id,
t.next_execute_time = date_add(now(),interval t.next_interval second)
where t.max_execute_time > now() 
and t.next_execute_time <= now() 
and t.status in(20,30)
and t.count <= 5
limit 1000`

const SQLQueryWaitProcess = `
select t.queue_name,t.msg_content content 
from tsk_system_task t
where t.batch_id=@batch_id 
and t.next_execute_time > now()`

const SQLClearTask = `delete from tsk_system_task where delete_time < now() and status in (0, 90)`

const SQLClearSEQ = `delete from tsk_system_seq
where 1=1 &seq_id`

const SQLFailedTask = `
UPDATE tsk_system_task t SET 
t.delete_time = DATE_ADD(NOW(),INTERVAL CASE WHEN t.delete_interval=0 THEN 604800 ELSE t.delete_interval END SECOND),
t.status = 90
WHERE t.max_execute_time > DATE_SUB(NOW(),INTERVAL 2 DAY)
AND (t.max_execute_time < NOW() OR t.count > 5) 
AND t.status IN (20, 30)
LIMIT 1000
`
