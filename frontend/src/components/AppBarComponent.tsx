import { AppBar, Container, Toolbar, IconButton, Typography, Avatar, Menu, MenuItem, Divider, CircularProgress } from '@mui/material';
import { WaterDrop as WatermarkIcon, AccountCircle } from '@mui/icons-material';
import { useState, useEffect } from 'react';
import theme from '../theme';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';

const AppBarComponent = () => {
  const navigate = useNavigate();
  const { user, logout } = useAuth();
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    setIsLoading(false);
  }, [user]);

  console.log('user from app bar: ', user);

  const handleMenu = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const handleSignIn = () => {
    handleClose();
    navigate('/signin');
  };

  const handleSignUp = () => {
    handleClose();
    navigate('/signup');
  };

  const handleLogout = () => {
    handleClose();
    logout();
    navigate('/');
  };

  const menuItems = user ? [
    <MenuItem key="username" disabled>{user.username}</MenuItem>,
    <Divider key="divider" />,
    <MenuItem key="account" onClick={handleClose}>Account</MenuItem>,
    <MenuItem key="logout" onClick={handleLogout}>Logout</MenuItem>
  ] : [
    <MenuItem key="signin" onClick={handleSignIn}>Sign In</MenuItem>,
    <MenuItem key="signup" onClick={handleSignUp}>Create Account</MenuItem>
  ];

  return (
    <AppBar position="static" color="primary" elevation={0}>
      <Container maxWidth="xl">
        <Toolbar disableGutters>
          <IconButton 
            edge="start" 
            color="inherit" 
            aria-label="menu" 
            sx={{ 
              mr: 2,
              color: theme.palette.secondary.main,
            }}
          >
            <WatermarkIcon fontSize="large" />
          </IconButton>
          <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
            Watermark Generator
          </Typography>
          <div>
            <IconButton onClick={handleMenu} sx={{ p: 0 }}>
              {isLoading ? (
                <CircularProgress size={24} color="inherit" />
              ) : user ? (
                <Avatar sx={{ bgcolor: theme.palette.secondary.main }}>
                  {user.username?.[0]}
                </Avatar>
              ) : (
                <AccountCircle fontSize="large" />
              )}
            </IconButton>
            {/* <IconButton onClick={handleMenu} sx={{ p: 0 }}>
              {user === undefined ? (
                <CircularProgress size={24} color="inherit" />
              ) : user ? (
                <Avatar sx={{ bgcolor: theme.palette.secondary.main }}>
                  {user.username?.[0]}
                </Avatar>
              ) : (
                <AccountCircle fontSize="large" />
              )}
            </IconButton> */}
            <Menu
              anchorEl={anchorEl}
              open={Boolean(anchorEl)}
              onClose={handleClose}
              anchorOrigin={{
                vertical: 'bottom',
                horizontal: 'right',
              }}
              transformOrigin={{
                vertical: 'top',
                horizontal: 'right',
              }}
            >
              {menuItems}
            </Menu>
          </div>
        </Toolbar>
      </Container>
    </AppBar>
  );
};

export default AppBarComponent;