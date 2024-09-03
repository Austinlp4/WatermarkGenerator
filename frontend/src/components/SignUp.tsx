import React, { useState } from 'react';
import { Box, Button, TextField, Typography, Link, Grid } from '@mui/material';
import { keyframes } from '@emotion/react';
import { useNavigate } from 'react-router-dom';

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

const SignUp = () => {
  const navigate = useNavigate();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setError('');

    const payload = { email, password: password.trim() };
    console.log('Sending payload:', payload);

    try {
      const response = await fetch(`${import.meta.env.VITE_API_URL}/api/register`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(payload),
      });

      if (response.ok) {
        navigate('/signin');
      } else {
        const contentType = response.headers.get("content-type");
        if (contentType && contentType.indexOf("application/json") !== -1) {
          const data = await response.json();
          setError(data.message || 'Registration failed');
        } else {
          const text = await response.text();
          setError(text || 'Registration failed');
        }
      }
    } catch (error) {
      console.error('Fetch error:', error);
      setError('An unexpected error occurred');
    }
  };

  return (
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
            Join Watermark Wizard
          </Typography>
          <Typography variant="body1" sx={{ mb: 4, color: 'text.secondary' }}>
            Sign up to start your creative journey.
          </Typography>
          <Box component="form" onSubmit={handleSubmit}>
            <TextField
              label="Email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              fullWidth
              margin="normal"
            />
            <TextField
              label="Password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              fullWidth
              margin="normal"
            />
            {error && <Typography color="error" sx={{ mt: 2 }}>{JSON.parse(error)?.message}</Typography>}
            <Button type="submit" fullWidth variant="contained" sx={{ mt: 3, mb: 2, py: 1.5, bgcolor: 'primary.main', '&:hover': { bgcolor: 'primary.dark' } }}>
              Sign Up
            </Button>
            <Grid container justifyContent="center">
              <Grid item>
                <Link href="/signin" variant="body2" sx={{color: 'primary.main'}}>Already have an account? Sign In</Link>
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
          Create Waves with Your Art
        </Typography>
        <Typography variant="h6" sx={{ 
          mb: 4, 
          maxWidth: '600px', 
          zIndex: 1,
          fontStyle: 'italic',
          color: '#e6f7ff',
        }}>
          Join our community of creators and protect your visual masterpieces 
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
  );
};

export default SignUp;