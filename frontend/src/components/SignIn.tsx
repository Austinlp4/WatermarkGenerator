import React, { useState } from 'react';
import { Box, Button, TextField, Typography, Link, Grid } from '@mui/material';
import { keyframes } from '@emotion/react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';
import SEO from './SEO';

const waveAnimation = keyframes`
  0% {
    transform: translateX(0) translateZ(0);
  }
  50% {
    transform: translateX(-30%) translateZ(0);
  }
  100% {
    transform: translateX(0) translateZ(0);
  }
`;

const wavySvg = `
  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 3200 200">
    <path fill="#003399" d="M0,100 C800,200 1600,0 2400,100 C3200,200 3200,100 3200,100 L3200,200 L0,200 Z" />
  </svg>
`;

const encodedWavySvg = encodeURIComponent(wavySvg);

const SignIn = () => {
  const navigate = useNavigate();
  const { login } = useAuth();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setError('');
    try {
      const response = await fetch(`${import.meta.env.VITE_API_URL}/api/signin`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password: password.trim() }),
      });

      if (!response.ok) {
        const errorData = await response.json();
        setError(errorData.message || 'Sign in failed');
        return;
      }

      const data = await response.json();
      login(data);
      navigate('/');
    } catch (err) {
      console.error('Error during sign in:', err);
      setError('An error occurred during sign in');
    }
  };

  return (
    <>
      <SEO 
        title="Sign In"
        description="Sign in to your Watermark Generator account to protect your images and manage your watermarks."
        canonicalUrl="https://watermark-generator.com/signin"
      />
      <Grid container component="main" sx={{ height: '80vh', overflow: 'hidden' }}>
        <Grid item xs={12} md={6} sx={{
          display: 'flex',
          flexDirection: 'column',
          justifyContent: 'center',
          p: 4,
          bgcolor: 'background.paper',
        }}>
          <Box sx={{ maxWidth: 400, mx: 'auto' }}>
            <Typography variant="h4" component="h1" gutterBottom fontWeight="bold" sx={{ color: 'primary.main' }}>
              Welcome to Watermark Wizard
            </Typography>
            <Typography variant="body1" sx={{ mb: 4, color: 'text.secondary' }}>
              Sign in to dive into your creative flow.
            </Typography>
            <Box component="form" onSubmit={handleSubmit}>
              <TextField
                fullWidth
                label="Email"
                variant="outlined"
                margin="normal"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
              />
              <TextField
                fullWidth
                label="Password"
                type="password"
                variant="outlined"
                margin="normal"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
              />
              {error && <Typography color="error" sx={{ mt: 2 }}>{error}</Typography>}
              <Button 
                type="submit" 
                fullWidth 
                variant="contained" 
                sx={{ mt: 3, mb: 2, py: 1.5, bgcolor: 'primary.main', '&:hover': { bgcolor: 'primary.dark' } }}
              >
                Sign In
              </Button>
              <Grid container justifyContent="space-between">
                <Grid item>
                  <Link href="#" variant="body2" sx={{color: 'primary.main'}}>Forgot password?</Link>
                </Grid>
                <Grid item>
                  <Link href="/signup" variant="body2" sx={{color: 'primary.main'}}>Don't have an account? Sign Up</Link>
                </Grid>
              </Grid>
            </Box>
          </Box>
        </Grid>
        <Grid item xs={12} md={6} sx={{
          position: 'relative',
          bgcolor: '#0099ff',
          overflow: 'hidden',
          display: 'flex',
          flexDirection: 'column',
          justifyContent: 'center',
          alignItems: 'center',
          color: 'white',
          textAlign: 'center',
          p: 4,
        }}>
          <Typography variant="h2" component="h1" sx={{ 
            mb: 2, 
            fontWeight: 'bold', 
            zIndex: 1,
            textShadow: '2px 2px 4px rgba(0,0,0,0.3)',
          }}>
            Make Waves with Your Designs
          </Typography>
          <Typography variant="h6" sx={{ 
            mb: 4, 
            maxWidth: '600px', 
            zIndex: 1,
            fontStyle: 'italic',
            color: '#e6f7ff',
          }}>
            Splash your mark on every image. Protect your visual masterpieces 
            with ease, style, and a touch of magic.
          </Typography>
          <Box sx={{ 
            display: 'flex', 
            gap: 2, 
            zIndex: 1 
          }}>
            <Button 
              variant="outlined" 
              size="large" 
              sx={{ 
                color: 'white', 
                borderColor: 'white',
                '&:hover': { bgcolor: 'rgba(255,255,255,0.1)' },
              }}
            >
              Learn More
            </Button>
          </Box>
          
          <Box sx={{
            position: 'absolute',
            left: 0,
            right: 0,
            bottom: 0,
            height: '100%',
            background: 'linear-gradient(180deg, #0099ff 0%, #0066cc 100%)',
          }}>
            <Box sx={{
              position: 'absolute',
              left: -800,
              right: -800,
              bottom: 0,
              height: '50%',
              backgroundImage: `url("data:image/svg+xml,${encodedWavySvg}")`,
              backgroundRepeat: 'repeat-x',
              backgroundSize: '3200px 200px',
              backgroundPosition: '0 bottom',
              animation: `${waveAnimation} 50s linear infinite`,
            }} />
            <Box sx={{
              position: 'absolute',
              left: -800,
              right: -800,
              bottom: '10px',
              height: '50%',
              backgroundImage: `url("data:image/svg+xml,${encodedWavySvg}")`,
              backgroundRepeat: 'repeat-x',
              backgroundSize: '3200px 200px',
              backgroundPosition: '0 bottom',
              animation: `${waveAnimation} 30s linear infinite`,
              opacity: 0.5,
            }} />
          </Box>
        </Grid>
      </Grid>
    </>
  );
};

export default SignIn;