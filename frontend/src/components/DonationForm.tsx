import React, { useState, useCallback } from 'react';
import { 
  Box, 
  Button, 
  CircularProgress, 
  Dialog, 
  DialogActions, 
  DialogContent, 
  DialogTitle, 
  TextField, 
  Typography 
} from '@mui/material';
import { Favorite as FavoriteIcon } from '@mui/icons-material';
import axios from 'axios';
import { loadStripe } from '@stripe/stripe-js';

const DonationForm = () => {
  const [amount, setAmount] = useState(5.00);
  const [isProcessing, setIsProcessing] = useState(false);
  const [open, setOpen] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleOpen = () => setOpen(true);
  const handleClose = () => {
    setOpen(false);
    setError(null);
  };

  const handleSubmit = useCallback(async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setIsProcessing(true);
    setError(null);

    try {
      const response = await axios.post('/api/create-checkout-session', {
        amount: Math.round(amount * 100),
        currency: 'usd',
      });

      const { sessionId } = response.data;

      const stripe = await loadStripe(import.meta.env.VITE_STRIPE_PUBLISHABLE_KEY);
      if (stripe) {
        const { error } = await stripe.redirectToCheckout({ sessionId });
        if (error) {
          setError(error.message || 'An error occurred. Please try again.');
        }
      } else {
        setError('Failed to load Stripe.');
      }
    } catch (error) {
      console.error('Error creating checkout session:', error);
      setError('An unexpected error occurred. Please try again.');
    } finally {
      setIsProcessing(false);
    }
  }, [amount]);

  return (
    <>
      <Button 
        variant="contained" 
        color="secondary" 
        onClick={handleOpen}
        startIcon={<FavoriteIcon />}
        sx={{ my: 2, width: '100%' }}
      >
        Support Us
      </Button>
      <Dialog open={open} onClose={handleClose} maxWidth="sm" fullWidth>
        <DialogTitle>
          <Box display="flex" alignItems="center">
            <FavoriteIcon sx={{ mr: 1, color: 'secondary.main' }} />
            <Typography variant="h5">Make a Donation</Typography>
          </Box>
        </DialogTitle>
        <form onSubmit={handleSubmit}>
          <DialogContent>
            <TextField
              fullWidth
              label="Donation Amount ($)"
              type="number"
              value={amount.toFixed(2)}
              onChange={(e) => setAmount(Number(e.target.value))}
              InputProps={{ inputProps: { min: 0.01, step: 0.01 } }}
              sx={{ mb: 2 }}
            />
            <Typography variant="body2" color="textSecondary">
              Your donation will help us maintain and improve Watermark Wizard.
            </Typography>
            {error && (
              <Typography color="error" variant="body2" sx={{ mt: 2 }}>
                {error}
              </Typography>
            )}
          </DialogContent>
          <DialogActions>
            <Button onClick={handleClose} sx={{ color: 'text.primary' }}>
              Cancel
            </Button>
            <Button
              type="submit"
              variant="contained"
              color="primary"
              disabled={isProcessing || amount <= 0}
            >
              {isProcessing ? <CircularProgress size={24} /> : `Donate $${amount.toFixed(2)}`}
            </Button>
          </DialogActions>
        </form>
      </Dialog>
    </>
  );
};

export default DonationForm;