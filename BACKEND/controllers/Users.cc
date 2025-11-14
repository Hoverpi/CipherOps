#include "Users.h"

void Users::loginUser(const HttpRequestPtr &req, std::function<void (const HttpResponsePtr &)> &&callback,
               std::string &&userId, const std::string &password) {
    
}


void Users::registerUser(const HttpRequestPtr &req, std::function<void (const HttpResponsePtr &)> &&callback,
                std::string userId, const std::string &token) const {

}

                