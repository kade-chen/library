# 私钥（只给 Auth 服务）
openssl genrsa -out jwt_private.pem 2048

# 公钥（给所有服务）
openssl rsa -in jwt_private.pem -pubout -out jwt_public.pem

# 更新
domain修改成organization的域名