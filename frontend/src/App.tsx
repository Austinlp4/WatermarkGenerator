import React, { useState, useCallback, useEffect } from 'react';
import { Routes, Route } from 'react-router-dom';
import { 
  Box, 
  Container, 
  Typography, 
  TextField, 
  Button, 
  Slider, 
  Paper,
  Grid,
  Snackbar,
  Alert
} from '@mui/material';
import { CloudUpload, Download as DownloadIcon } from '@mui/icons-material';
import { ThemeProvider } from '@mui/material/styles';
import { loadStripe } from '@stripe/stripe-js';
import { Elements } from '@stripe/react-stripe-js';
import DonationForm from './components/DonationForm';
import SignIn from './components/SignIn';
import SignUp from './components/SignUp';
import theme from './theme';
import ImagePreview from './components/ImagePreview';
import AppBarComponent from './components/AppBarComponent';
import { useAuth } from './hooks/useAuth';
import AdminPage from './components/AdminPage';
import { useNavigate } from 'react-router-dom';

const stripePromise = loadStripe(import.meta.env.VITE_STRIPE_PUBLISHABLE_KEY);

function App() {
  const { user } = useAuth();
  const [file, setFile] = useState<File | null>(null);
  const [preview, setPreview] = useState<string | null>(null);
  const [watermarkText, setWatermarkText] = useState('');
  const [textColor, setTextColor] = useState('#ffffff');
  const [opacity, setOpacity] = useState(0.5);
  const [fontSize, setFontSize] = useState(100);
  const [spacing, setSpacing] = useState(200);
  const [isLoading, setIsLoading] = useState(false);
  const [watermarkedImage, setWatermarkedImage] = useState<string | null>(null);
  const [donationStatus, setDonationStatus] = useState<string | null>(null);
  const navigate = useNavigate();

  useEffect(() => {
    const urlParams = new URLSearchParams(window.location.search);
    const status = urlParams.get('donation');
    if (status) {
      setDonationStatus(status);
      window.history.replaceState({}, document.title, window.location.pathname);
    }
  }, []);

  const handleCloseSnackbar = () => {
    setDonationStatus(null);
  };

  const handleFileChange = useCallback((event: React.ChangeEvent<HTMLInputElement>) => {
    if (event.target.files && event.target.files[0]) {
      const selectedFile = event.target.files[0];
      setFile(selectedFile);
      const reader = new FileReader();
      reader.onloadend = () => {
        setPreview(reader.result as string);
      };
      reader.readAsDataURL(selectedFile);
    }
  }, []);

  const handleSubmit = useCallback(async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (!file || !user) return;

    setIsLoading(true);
    const formData = new FormData();
    formData.append('image', file);
    formData.append('text', watermarkText);
    formData.append('color', textColor);
    formData.append('opacity', opacity.toString());
    formData.append('fontSize', fontSize.toString());
    formData.append('spacing', spacing.toString());

    try {
      const response = await fetch('/api/watermark', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${user.token}`,
        },
        body: formData,
      });

      if (response.ok) {
        const blob = await response.blob();
        const url = URL.createObjectURL(blob);
        setPreview(url);
        setWatermarkedImage(url);
      } else {
        console.error('Error applying watermark');
      }
    } catch (error) {
      console.error('Error:', error);
    } finally {
      setIsLoading(false);
    }
  }, [file, watermarkText, textColor, opacity, fontSize, spacing, user]);

  const handleDownload = useCallback(() => {
    if (watermarkedImage && user) {
      window.location.href = `/api/download?path=${watermarkedImage}&token=${user.token}`;
    } else if (!user) {
      navigate('/signin');
    }
  }, [watermarkedImage, user, navigate]);

  return (
    <ThemeProvider theme={theme}>
        <Elements stripe={stripePromise}>
          <Box sx={{ display: 'flex', flexDirection: 'column', minHeight: '100vh', bgcolor: 'background.default' }}>
            <AppBarComponent />
            <Container maxWidth="xl" sx={{ mt: 4, mb: 4, flexGrow: 1, display: 'flex', flexDirection: 'column' }}>
              <Routes>
                <Route path="/" element={
                  <Grid container spacing={3}>
                    <Grid item xs={12} md={4}>
                      <Grid container direction="column" spacing={3}>
                        <Grid item>
                          <Paper elevation={3} sx={{ p: 3 }}>
                            <Typography variant="h6" gutterBottom color="secondary">Watermark Controls</Typography>
                            <Box component="form" onSubmit={handleSubmit} sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
                              <Button variant="outlined" component="label" startIcon={<CloudUpload />}>
                                Upload Image
                                <input type="file" hidden onChange={handleFileChange} accept="image/*" />
                              </Button>
                              
                              <TextField
                                fullWidth
                                label="Watermark Text"
                                value={watermarkText}
                                onChange={(e) => setWatermarkText(e.target.value)}
                              />
                              
                              <TextField
                                fullWidth
                                label="Text Color"
                                type="color"
                                value={textColor}
                                onChange={(e) => setTextColor(e.target.value)}
                              />
                              
                              <Box>
                                <Typography gutterBottom>Opacity</Typography>
                                <Slider
                                  value={opacity}
                                  onChange={(_, newValue) => setOpacity(newValue as number)}
                                  min={0}
                                  max={1}
                                  step={0.1}
                                />
                              </Box>
                              
                              <Box>
                                <Typography gutterBottom>Font Size</Typography>
                                <Slider
                                  value={fontSize}
                                  onChange={(_, newValue) => setFontSize(newValue as number)}
                                  min={10}
                                  max={200}
                                  step={1}
                                />
                              </Box>
                              
                              <Box>
                                <Typography gutterBottom>Spacing</Typography>
                                <Slider
                                  value={spacing}
                                  onChange={(_, newValue) => setSpacing(newValue as number)}
                                  min={50}
                                  max={800}
                                  step={5}
                                />
                              </Box>
                              
                              <Button type="submit" variant="contained" color="primary" disabled={!file}>
                                Apply Watermark
                              </Button>
                            </Box>
                          </Paper>
                        </Grid>
                        
                        <Grid item>
                          <Paper elevation={3} sx={{ p: 3 }}>
                            <Typography variant="h6" gutterBottom>Support the Project</Typography>
                            <DonationForm />
                          </Paper>
                        </Grid>
                      </Grid>
                    </Grid>
                    
                    <Grid item xs={12} md={8}>
                      <Paper elevation={3} sx={{ p: 3, height: '100%', display: 'flex', flexDirection: 'column' }}>
                        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
                          <Typography variant="h6" gutterBottom color="secondary">Preview</Typography>
                          {watermarkedImage && (
                            <Button
                              variant="contained"
                              color="secondary"
                              startIcon={<DownloadIcon />}
                              onClick={handleDownload}
                              size="small"
                            >
                              Download Watermark
                            </Button>
                          )}
                        </Box>
                        <Box sx={{ flexGrow: 1, display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '400px' }}>
                          <ImagePreview
                            preview={preview}
                            isLoading={isLoading}
                            watermarkedImage={watermarkedImage}
                          />
                        </Box>
                      </Paper>
                    </Grid>
                  </Grid>
                } />
                <Route path="/signin" element={<SignIn />} />
                <Route path="/signup" element={<SignUp />} />
                <Route path="/admin" element={<AdminPage />} />
              </Routes>
            </Container>
            <Box
              component="footer"
              sx={{
                py: 3,
                px: 2,
                mt: 'auto',
                backgroundColor: theme.palette.background.default,
                color: theme.palette.text.primary,
              }}
            >
              <Container maxWidth="sm">
                <Typography variant="body2" color="text.secondary" align="center">
                  {'Copyright © '}
                    Watermark Generator{' '}
                  {new Date().getFullYear()}
                  {'.'}
                </Typography>
              </Container>
            </Box>
          </Box>
          
          {/* Success Snackbar */}
          <Snackbar
            open={donationStatus === 'success'}
            autoHideDuration={6000}
            onClose={handleCloseSnackbar}
            anchorOrigin={{ vertical: 'top', horizontal: 'center' }}
          >
            <Alert onClose={handleCloseSnackbar} severity="success" sx={{ width: '100%' }}>
              Thank you for your donation!
            </Alert>
          </Snackbar>

          {/* Cancelled Snackbar */}
          <Snackbar
            open={donationStatus === 'cancelled'}
            autoHideDuration={6000}
            onClose={handleCloseSnackbar}
            anchorOrigin={{ vertical: 'top', horizontal: 'center' }}
          >
            <Alert onClose={handleCloseSnackbar} severity="info" sx={{ width: '100%' }}>
              Donation cancelled. No charges were made.
            </Alert>
          </Snackbar>
        </Elements>
    </ThemeProvider>
  );
}

export default App;