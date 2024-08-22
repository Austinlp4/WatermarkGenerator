import React from 'react';
import { Dialog, DialogTitle, DialogContent, DialogActions, Button, Typography } from '@mui/material';
import { useNavigate } from 'react-router-dom';
import { useTheme } from '@mui/material/styles';

interface AuthModalProps {
  open: boolean;
  onClose: () => void;
}

const AuthModal: React.FC<AuthModalProps> = ({ open, onClose }) => {
  const navigate = useNavigate();
  const theme = useTheme();

  const handleSignIn = () => {
    onClose();
    navigate('/signin');
  };

  const handleSignUp = () => {
    onClose();
    navigate('/signup');
  };

  return (
    <Dialog open={open} onClose={onClose}>
      <DialogTitle>Authentication Required</DialogTitle>
      <DialogContent>
        <Typography>
          You need to be signed in to apply a watermark. Please sign in or create an account to continue.
        </Typography>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose} sx={{ color: theme.palette.text.secondary }}>
          Cancel
        </Button>
        <Button onClick={handleSignIn} variant="contained" color="primary">
          Sign In
        </Button>
        <Button onClick={handleSignUp} variant="outlined" color="secondary">
          Sign Up
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default AuthModal;