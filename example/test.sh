#!/bin/env bash

# 列举所有数据库
curl http://127.0.0.1:8000/schemas
# 统计所有数据库数量
curl http://127.0.0.1:8000/schemas?_count=*
# 查询hsb_youpin_mall 数据库表t_user所有数据
curl  http://127.0.0.1:8000/hsb_youpin_mall/SCHEMA/t_user

#分页查询
curl http://127.0.0.1:8000/hsb_youpin_mall/SCHEMA/t_user?_page=2&_page_size=10
# 条件查询(Fuser_id 值为 604337144)
curl http://127.0.0.1:8000/hsb_youpin_mall/SCHEMA/t_user?Fuser_id=604337144

# 插入记录,返回传递参数作为查询条件的记录列表（如果传递参数在表中唯一，则相当于返回新增记录）
curl -XPOST -d '{"Fuser_id":"123455","Fsource":"new test"}' http://127.0.0.1:8000/hsb_youpin_mall/SCHEMA/t_utm
# 更新记录
curl -XPUT -d '{"Fuser_id":"123455","Fsource":"new test9999"}' http://127.0.0.1:8000/hsb_youpin_mall/SCHEMA/t_utm?Fid=8
curl -XPATCH -d '{"Fuser_id":"123455","Fsource":"new test9999"}' http://127.0.0.1:8000/hsb_youpin_mall/SCHEMA/t_utm?Fid=8

# 删除记录 
curl -XDELETE http://127.0.0.1:8000/hsb_youpin_mall/SCHEMA/t_utm?Fid=7

# 使用sql模版
curl http://127.0.0.1:8000/_QUERIES/user/user_get?source1=test&source2=ee