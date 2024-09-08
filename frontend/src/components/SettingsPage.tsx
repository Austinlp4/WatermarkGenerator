import React, { useState, useEffect } from 'react';
import { 
  Box, 
  Typography, 
  Button, 
  Paper, 
  Grid, 
  Divider, 
  Dialog, 
  DialogActions, 
  DialogContent, 
  DialogContentText, 
  DialogTitle 
} from '@mui/material';
import { useAuth } from '../hooks/useAuth';
import { useNavigate } from 'react-router-dom';

const SettingsPage: React.FC = () => {
  const { user, refreshUser } = useAuth();
  const navigate = useNavigate();
  const [openDialog, setOpenDialog] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    if (!user) {
      navigate('/signin');
    }
  }, [user, navigate]);

  const handleCancelSubscription = async () => {
    setIsLoading(true);
    try {
      const response = await fetch(`${import.meta.env.VITE_API_URL}/api/cancel-subscription`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${user?.token}`,
        },
        body: JSON.stringify({ userId: user?.id }),
      });

      if (response.ok) {
        await refreshUser();
        setOpenDialog(false);
      } else {
        console.error('Failed to cancel subscription');
      }
    } catch (error) {
      console.error('Error cancelling subscription:', error);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Box sx={{ flexGrow: 1, padding: 3 }}>
      <Typography variant="h4" gutterBottom sx={{ color: 'primary.main' }}>
        Account Settings
      </Typography>
      <Paper elevation={3} sx={{ padding: 3, marginTop: 3 }}>
        <Grid container spacing={3}>
          <Grid item xs={12}>
            <Typography variant="h6">Personal Information</Typography>
            <Divider sx={{ my: 2 }} />
            <Typography><strong>Email:</strong> {user?.email}</Typography>
          </Grid>
          <Grid item xs={12}>
            <Typography variant="h6">Subscription Details</Typography>
            <Divider sx={{ my: 2 }} />
            <Typography><strong>Status:</strong> {user?.subscriptionStatus}</Typography>
            {user?.subscriptionStatus === 'active' && (
              <>
                <Typography><strong>Expires At:</strong> {new Date(user.subscriptionExpiresAt).toLocaleString()}</Typography>
                <Button 
                  variant="contained" 
                  color="secondary" 
                  onClick={() => setOpenDialog(true)}
                  sx={{ mt: 2 }}
                >
                  Cancel Subscription
                </Button>
              </>
            )}
            {user?.subscriptionStatus !== 'active' && (
              <Button 
                variant="contained" 
                color="primary" 
                onClick={() => navigate('/subscription')}
                sx={{ mt: 2 }}
              >
                Upgrade to Pro
              </Button>
            )}
          </Grid>
        </Grid>
      </Paper>

      <Dialog
        open={openDialog}
        onClose={() => setOpenDialog(false)}
      >
        <DialogTitle>Cancel Subscription</DialogTitle>
        <DialogContent>
          <DialogContentText>
            Are you sure you want to cancel your subscription? You will lose access to premium features at the end of your current billing cycle.
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenDialog(false)} color="primary">
            No, Keep My Subscription
          </Button>
          <Button onClick={handleCancelSubscription} color="secondary" disabled={isLoading}>
            {isLoading ? 'Cancelling...' : 'Yes, Cancel Subscription'}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default SettingsPage;