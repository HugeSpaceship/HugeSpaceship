--name GetSession: one
SELECT * FROM sessions WHERE id = sqlc.arg(id);

--name GetSessionFromToken: one
SELECT * FROM sessions WHERE token = sqlc.arg(token);

--name GetUserSessions: many
SELECT * FROM sessions WHERE userid = sqlc.arg(userID);

--name NewSession
INSERT INTO sessions(userid, ip, token, game, platform, expiry) VALUES (
                                                                        sqlc.arg(userID),
                                                                        sqlc.arg(ip),
                                                                        sqlc.arg(token),
                                                                            sqlc.arg(game),
                                                                            sqlc.arg(platform),
                                                                        sqlc.arg(expiry)
                                                                       );