import { useEffect } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { Box, Typography, Button, CircularProgress } from '@mui/material';
import CheckCircleOutlineIcon from '@mui/icons-material/CheckCircleOutline';
import { useAuth } from '../hooks/useAuth';

const SubscriptionSuccess = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const { user, refreshUser } = useAuth();

  useEffect(() => {
    const queryParams = new URLSearchParams(location.search);
    const sessionId = queryParams.get('session_id');

    if (sessionId) {
      // You might want to verify the session on your backend here
      refreshUser();
    }
  }, [location, refreshUser]);

  if (!user) {
    return <CircularProgress />;
  }

  return (
    <Box
      display="flex"
      flexDirection="column"
      alignItems="center"
      justifyContent="center"
      minHeight="80vh"
    >
      <CheckCircleOutlineIcon color="primary" style={{ fontSize: 80 }} />
      <Typography variant="h4" gutterBottom sx={{ color: 'primary.main' }}>
        Subscription Successful!
      </Typography>
      <Typography variant="body1" align="center" paragraph sx={{ color: 'text.secondary' }}>
        Thank you for subscribing to our Pro plan. You now have access to all premium features.
      </Typography>
      <Button variant="contained" color="primary" onClick={() => navigate('/')}>
        Go to Dashboard
      </Button>
    </Box>
  );
};

export default SubscriptionSuccess;