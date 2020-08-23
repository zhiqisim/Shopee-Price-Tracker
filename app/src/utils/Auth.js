import React, { useEffect, useState } from "react";
import API from "../utils/API";

export const AuthContext = React.createContext();

export const AuthProvider = ({ children }) => {
  const [currentUser, setCurrentUser] = useState(null);

  const checkForAuth = () => {
    API.get('/user/is-auth', { withCredentials: true })
        .then(response => {
            console.log(response.data);
            if (response.data.message === "success") {
                setCurrentUser(true)
            }
        })
        .catch(error => {
            console.log(error);
        });
}

useEffect(() => {
    checkForAuth();
},);

  return (
    <AuthContext.Provider
      value={{
        currentUser,
        setCurrentUser
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};
