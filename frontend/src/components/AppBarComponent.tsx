import { AppBar, Container, Toolbar, IconButton, Typography, Avatar, Menu, MenuItem, Divider, CircularProgress, Button, Chip } from '@mui/material';
import { WaterDrop as WatermarkIcon, AccountCircle, Star as StarIcon } from '@mui/icons-material';
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

  const handleSubscription = () => {
    navigate('/subscription');
  };

  const menuItems = user ? [
    <Divider key="divider" />,
    <MenuItem key="account" onClick={handleClose}>Account</MenuItem>,
    <MenuItem key="settings" onClick={() => { handleClose(); navigate('/settings'); }}>Settings</MenuItem>,
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
            onClick={() => navigate('/')}
          >
            <WatermarkIcon fontSize="large" />
          </IconButton>
          <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
            Watermark Generator
          </Typography>
          {user && user.subscriptionStatus === 'active' ? (
            <Chip
              icon={<StarIcon fontSize="small" />}
              label="PRO"
              sx={{
                mr: 2,
                background: 'linear-gradient(45deg, #3a86ff 30%, #023e8a 90%)',
                color: 'white',
                fontWeight: 'bold',
                '& .MuiChip-icon': {
                  color: 'white',
                },
              }}
            />
          ) : (
            <Button
              variant="contained"
              onClick={handleSubscription}
              sx={{
                mr: 2,
                fontWeight: 'bold',
                borderRadius: '50px',
                padding: '6px 16px',
                background: 'linear-gradient(45deg, #3a86ff 30%, #023e8a 90%)',
                color: 'white',
                textTransform: 'none',
                fontSize: '0.875rem',
                border: 'none',
                boxShadow: 'none',
              }}
            >
              Upgrade to Pro
            </Button>
          )}
          <IconButton onClick={handleMenu} sx={{ p: 0 }}>
            {isLoading ? (
              <CircularProgress size={24} color="inherit" />
            ) : user ? (
              <Avatar sx={{ bgcolor: '#7fffd4', color: '#1a202c', fontWeight: 'bold' }}>
                {user.email?.[0].toUpperCase()}
              </Avatar>
            ) : (
              <AccountCircle fontSize="large" />
            )}
          </IconButton>
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
            PaperProps={{
              style: {
                maxWidth: '300px', // Adjust this value as needed
                width: '100%',
              },
            }}
          >
            {menuItems}
          </Menu>
        </Toolbar>
      </Container>
    </AppBar>
  );
};

export default AppBarComponent;