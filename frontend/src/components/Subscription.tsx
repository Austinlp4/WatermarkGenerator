import { Box, Button, Card, CardContent as MuiCardContent, Container, Grid, Typography, styled } from '@mui/material';
import { motion } from 'framer-motion';
import CheckCircleOutlineIcon from '@mui/icons-material/CheckCircleOutline';
import StarIcon from '@mui/icons-material/Star';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import { loadStripe } from '@stripe/stripe-js';
import SEO from './SEO';

const AnimatedCard = styled(motion(Card))(() => ({
  height: '100%',
  display: 'flex',
  flexDirection: 'column',
  borderRadius: '16px',
  boxShadow: '0 4px 30px rgba(0, 0, 0, 0.1)',
  backdropFilter: 'blur(5px)',
  backgroundColor: 'rgba(255, 255, 255, 0.1)',
  border: '1px solid rgba(255, 255, 255, 0.3)',
  overflow: 'hidden',
  transition: 'all 0.3s ease-in-out',
  '&:hover': {
    transform: 'translateY(-10px)',
    boxShadow: '0 20px 30px rgba(0, 0, 0, 0.15)',
  },
}));

const CardHeader = styled(Box)(({ theme }) => ({
  padding: theme.spacing(3),
  background: 'linear-gradient(45deg, #2196F3 30%, #21CBF3 90%)',
  color: 'white',
}));

const FeatureItem = styled(Box)(({ theme }) => ({
  display: 'flex',
  alignItems: 'center',
  marginBottom: theme.spacing(1),
}));

const AnimatedButton = styled(motion(Button))({
  marginTop: '16px',
  padding: '12px',
  fontWeight: 'bold',
  borderRadius: '8px',
});

const StyledCardContent = styled(MuiCardContent)(({ theme }) => ({
  flexGrow: 1,
  display: 'flex',
  flexDirection: 'column',
  justifyContent: 'space-between',
  padding: theme.spacing(3),
}));

const FeatureList = styled(Box)(({ theme }) => ({
  flexGrow: 1,
  marginBottom: theme.spacing(2),
}));

const Subscription = () => {
  const navigate = useNavigate();
  const { user } = useAuth() || {};

  const handleSubscribe = async (tier: 'free' | 'pro') => {
    if (tier === 'free') {
      navigate(user ? '/' : '/login');
    } else {
      try {
        if (!user || !user.id) {
          throw new Error('User not authenticated');
        }

        console.log('Sending subscription request...');
        const response = await fetch(`${import.meta.env.VITE_API_URL}/api/create-subscription`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${user.token}`,
          },
          body: JSON.stringify({ userId: user.id }),
        });

        console.log('Response status:', response.status);
        const responseText = await response.text();
        console.log('Response text:', responseText);

        if (!response.ok) {
          throw new Error(`Failed to create subscription: ${responseText}`);
        }

        let sessionData;
        try {
          sessionData = JSON.parse(responseText);
        } catch (parseError) {
          console.error('Error parsing response:', parseError);
          throw new Error('Invalid response from server');
        }

        console.log('Session data:', sessionData);

        if (!sessionData.sessionId) {
          throw new Error('No sessionId in response');
        }

        // Redirect to Stripe Checkout
        const stripe = await loadStripe(import.meta.env.VITE_STRIPE_PUBLISHABLE_KEY);
        if (stripe) {
          await stripe.redirectToCheckout({ sessionId: sessionData.sessionId });
        } else {
          throw new Error('Failed to load Stripe');
        }
      } catch (error) {
        console.error('Error in subscription process:', error);
        // Handle error (e.g., show error message to user)
      }
    }
  };

  return (
    <>
      <SEO 
        title="Subscription Plans"
        description="Choose the perfect Watermark Generator subscription plan for your needs. Protect your images with our free and pro options."
        canonicalUrl="https://watermark-generator.com/subscription"
      />
      <Container maxWidth="lg" sx={{ mt: 8, mb: 8 }}>
        <Typography variant="h3" align="center" gutterBottom component={motion.h3}
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5 }}
          sx={{ fontWeight: 'bold', mb: 4 }}
        >
          Choose Your Perfect Plan
        </Typography>
        <Grid container spacing={4} justifyContent="center">
          {/* Free Tier Card */}
          <Grid item xs={12} md={6}>
            <AnimatedCard initial={{ opacity: 0, x: -50 }} animate={{ opacity: 1, x: 0 }} transition={{ duration: 0.5 }}>
              <CardHeader>
                <Typography variant="h4" align="center" gutterBottom>
                  Free Tier
                </Typography>
                <Typography variant="subtitle1" align="center">
                  Perfect for getting started
                </Typography>
              </CardHeader>
              <StyledCardContent>
                <div style={{ display: 'flex', flexDirection: 'column', justifyContent: 'space-between', height: '100%' }}>
                  <FeatureList>
                    <FeatureItem>
                      <CheckCircleOutlineIcon color="primary" sx={{ mr: 1 }} />
                      <Typography variant="body1">1 high quality watermark generation daily</Typography>
                    </FeatureItem>
                    <FeatureItem>
                      <CheckCircleOutlineIcon color="primary" sx={{ mr: 1 }} />
                      <Typography variant="body1">Customizable text or image watermarks</Typography>
                    </FeatureItem>
                    <FeatureItem>
                      <CheckCircleOutlineIcon color="primary" sx={{ mr: 1 }} />
                      <Typography variant="body1">Instant delivery</Typography>
                    </FeatureItem>
                  </FeatureList>
                  <Typography variant="body2" color="textSecondary" sx={{ mb: 2, mt: 'auto' }}>
                    Start protecting your images today at no cost!
                  </Typography>
                </div>
                <AnimatedButton
                  fullWidth
                  variant="outlined"
                  color="primary"
                  onClick={() => handleSubscribe('free')}
                  whileHover={{ scale: 1.05 }}
                  whileTap={{ scale: 0.95 }}
                >
                  Get Started for Free
                </AnimatedButton>
              </StyledCardContent>
            </AnimatedCard>
          </Grid>
          {/* Pro Tier Card */}
          <Grid item xs={12} md={6}>
            <AnimatedCard initial={{ opacity: 0, x: 50 }} animate={{ opacity: 1, x: 0 }} transition={{ duration: 0.5 }}>
              <CardHeader sx={{ 
                background: 'linear-gradient(45deg, #64FFDA 30%, #00B8D4 90%)',
                color: '#1A202C' // Dark color for better contrast
              }}>
                <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', mb: 1 }}>
                  <StarIcon sx={{ mr: 1 }} />
                  <Typography variant="h4" align="center">
                    Pro Tier
                  </Typography>
                </Box>
                <Typography variant="subtitle1" align="center">
                  Unlock unlimited potential
                </Typography>
              </CardHeader>
              <StyledCardContent>
                <Typography variant="h5" align="center" gutterBottom color="secondary" sx={{ fontWeight: 'bold' }}>
                  $4.99/month
                </Typography>
                <Box sx={{ mb: 3 }}>
                  <FeatureItem>
                    <CheckCircleOutlineIcon color="secondary" sx={{ mr: 1 }} />
                    <Typography variant="body1">All Free Tier features</Typography>
                  </FeatureItem>
                  <FeatureItem>
                    <CheckCircleOutlineIcon color="secondary" sx={{ mr: 1 }} />
                    <Typography variant="body1">Bulk watermark generation</Typography>
                  </FeatureItem>
                  <FeatureItem>
                    <CheckCircleOutlineIcon color="secondary" sx={{ mr: 1 }} />
                    <Typography variant="body1">Unlimited generations a day</Typography>
                  </FeatureItem>
                  <FeatureItem>
                    <CheckCircleOutlineIcon color="secondary" sx={{ mr: 1 }} />
                    <Typography variant="body1">Advanced watermark placement and options</Typography>
                  </FeatureItem>
                </Box>
                <Typography variant="body2" color="textSecondary" sx={{ mb: 2 }}>
                  Take your image protection to the next level!
                </Typography>
                <AnimatedButton
                  fullWidth
                  variant="contained"
                  sx={{
                    background: 'linear-gradient(45deg, #64FFDA 30%, #00B8D4 90%)',
                    color: '#1A202C', // Dark color to match the header text
                    '&:hover': {
                      background: 'linear-gradient(45deg, #00B8D4 30%, #64FFDA 90%)',
                    }
                  }}
                  onClick={() => handleSubscribe('pro')}
                  whileHover={{ scale: 1.05 }}
                  whileTap={{ scale: 0.95 }}
                >
                  Upgrade to Pro
                </AnimatedButton>
              </StyledCardContent>
            </AnimatedCard>
          </Grid>
        </Grid>
      </Container>
    </>
  );
};

export default Subscription;