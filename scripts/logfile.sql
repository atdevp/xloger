CREATE TABLE `logfile_task` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `logset_id` varchar(36) NOT NULL,
  `logset_name` varchar(255) NOT NULL,
  `logset_split` varchar(10) NOT NULL,
  `agent_ip` varchar(15) NOT NULL,
  `agent_role` varchar(10) NOT NULL,
  `schedue_time` varchar(200) NOT NULL,
  `temp_path` varchar(2000) NOT NULL,
  `hdfs_path` varchar(2000) NOT NULL,
  `status` int(11) NOT NULL,
  `utime` datetime(6) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `logfile_task_logset_id_c5cbf3ce_uniq` (`logset_id`,`agent_ip`,`agent_role`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;



CREATE TABLE `logfile_host` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `logset_id` varchar(36) NOT NULL,
  `agent_role` varchar(10) NOT NULL,
  `host` varchar(15) NOT NULL,
  `passwd` varchar(15) NOT NULL,
  `file_path` varchar(2000) NOT NULL,
  `pl_stime` datetime(6) NOT NULL,
  `pl_etime` datetime(6) NOT NULL,
  `pl_state` varchar(10) NOT NULL,
  `pu_stime` datetime(6) NOT NULL,
  `pu_etime` datetime(6) NOT NULL,
  `pu_state` varchar(10) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
