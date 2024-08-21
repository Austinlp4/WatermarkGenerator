import React, { createContext, useState, ReactNode, useEffect } from 'react';

interface User {
  username: string;
  token: string;
}

interface AuthContextType {
  user: User | null;
  login: (user: User) => void;
  logout: () => void;
}

export const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchUser = async () => {
      const token = localStorage.getItem('token');
      console.log('Initial token from localStorage:', token);

      if (token) {
        try {
          console.log('Attempting to fetch user with token:', token);
          const response = await fetch('http://localhost:8080/api/current-user', {
            headers: {
              'Authorization': `Bearer ${token}`,
            },
          });

          if (response.ok) {
            const responseText = await response.text();
            if (responseText) {
              const userData = JSON.parse(responseText);
              console.log('User data fetched successfully:', userData);
              setUser({ ...userData, token });
            } else {
              console.error('Empty response from server');
              setError('Empty response from server');
              localStorage.removeItem('token');
            }
          } else {
            console.error('Unexpected response status:', response.status);
            const responseText = await response.text();
            console.error('Response text:', responseText);
            setError(`Failed to fetch user: ${response.status} ${responseText}`);
            localStorage.removeItem('token');
          }
        } catch (error) {
          console.error('Failed to fetch user:', error);
          setError(`Network error: ${(error as Error).message || 'Unknown error'}`);
          localStorage.removeItem('token');
        }
      } else {
        console.log('No token found in localStorage');
      }
    };

    fetchUser();
  }, []);

  const login = (userData: User) => {
    console.log('Logging in user:', userData);
    setUser(userData);
    localStorage.setItem('token', userData.token);
    console.log('Token set in localStorage:', userData.token);
  };

  const logout = () => {
    console.log('Logging out user');
    setUser(null);
    localStorage.removeItem('token');
    console.log('Token removed from localStorage');
  };

  if (error) {
    return (
      <div>
        <p>Error: {error}</p>
        <button onClick={() => { setError(null); logout(); }}>Clear error and logout</button>
      </div>
    );
  }

  return (
    <AuthContext.Provider value={{ user, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
};