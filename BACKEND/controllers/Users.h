#pragma once

#include <drogon/HttpController.h>

using namespace drogon;

class Users : public drogon::HttpController<Users>
{
public:
    METHOD_LIST_BEGIN
    // use METHOD_ADD to add your custom processing function here;
    METHOD_ADD(Users::loginUser,"/login",Post);
    METHOD_ADD(Users::registerUser,"/register",Post);
    METHOD_ADD(User::getInfo,"/{1}/info?token={2}",Get);
    METHOD_LIST_END
    // your declaration of processing function maybe like this:
    void loginUser(const HttpRequestPtr &req, std::function<void (const HttpResponsePtr &)> &&callback,
               std::string &&userId, const std::string &password);
    void registerUser(const HttpRequestPtr &req, std::function<void (const HttpResponsePtr &)> &&callback,
                 std::string userId, const std::string &token) const;
    void getInfo(const HttpRequestPtr &req, std::function<void (const HttpResponsePtr &)> &&callback,
                   std::string userId, const std::string &token) const;
};
