import { Box, Typography, Button } from '@mui/material';
import CancelOutlinedIcon from '@mui/icons-material/CancelOutlined';
import { useNavigate } from 'react-router-dom';

const SubscriptionCancel = () => {
  const navigate = useNavigate();

  return (
    <Box
      display="flex"
      flexDirection="column"
      alignItems="center"
      justifyContent="center"
      minHeight="80vh"
    >
      <CancelOutlinedIcon color="error" style={{ fontSize: 80 }} />
      <Typography variant="h4" gutterBottom>
        Subscription Cancelled
      </Typography>
      <Typography variant="body1" align="center" paragraph>
        Your subscription process was cancelled. If you have any questions or concerns, please contact our support team.
      </Typography>
      <Button variant="contained" color="primary" onClick={() => navigate('/subscription')}>
        Try Again
      </Button>
    </Box>
  );
};

export default SubscriptionCancel;