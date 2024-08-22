import { useEffect, useState } from 'react';
import { Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Paper, Typography } from '@mui/material';

interface User {
  id: string;
  username: string;
  firstName: string;
  email: string;
}

const AdminPage = () => {
  const [users, setUsers] = useState<User[]>([]);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchUsers = async () => {
      try {
        const response = await fetch(`${import.meta.env.VITE_API_URL}/api/users`, {
          credentials: 'include',
        });
        if (response.ok) {
          const usersData = await response.json();
          setUsers(usersData);
        } else {
          setError('Failed to fetch users');
        }
      } catch (error) {
        console.error(error)
        setError('An error occurred while fetching users');
      }
    };

    fetchUsers();
  }, []);

  return (
    <TableContainer component={Paper}>
      <Typography variant="h4" component="h1" gutterBottom>
        Users
      </Typography>
      {error && <Typography color="error">{error}</Typography>}
      <Table>
        <TableHead>
          <TableRow>
            <TableCell>ID</TableCell>
            <TableCell>Username</TableCell>
            <TableCell>First Name</TableCell>
            <TableCell>Email</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {users.map((user) => (
            <TableRow key={user.id}>
              <TableCell>{user.id}</TableCell>
              <TableCell>{user.username}</TableCell>
              <TableCell>{user.firstName}</TableCell>
              <TableCell>{user.email}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </TableContainer>
  );
};

export default AdminPage;