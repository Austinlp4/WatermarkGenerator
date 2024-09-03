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
  Alert,
  Tabs,
  Tab,
  Checkbox,
  FormControlLabel,
  Chip
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
import AuthModal from './components/AuthModal';
import Subscription from './components/Subscription';
import SubscriptionSuccess from './components/SubscriptionSuccess';
import SubscriptionCancel from './components/SubscriptionCancel';
import SettingsPage from './components/SettingsPage';
import BackgroundBlobs from './components/BackgroundBlobs';

const stripePromise = loadStripe(import.meta.env.VITE_STRIPE_PUBLISHABLE_KEY);

function App() {
  const { user } = useAuth();
  const [file, setFile] = useState<File | null>(null);
  const [watermarkText, setWatermarkText] = useState('');
  const [textColor, setTextColor] = useState('#ffffff');
  const [opacity, setOpacity] = useState(0.5);
  const [fontSize, setFontSize] = useState(100);
  const [spacing, setSpacing] = useState(500);
  const [isLoading, setIsLoading] = useState(false);
  const [watermarkedImage, setWatermarkedImage] = useState<string | null>(null);
  const [donationStatus, setDonationStatus] = useState<string | null>(null);
  const [showAuthModal, setShowAuthModal] = useState(false);
  const [watermarkImage, setWatermarkImage] = useState<File | null>(null);
  const [watermarkSize, setWatermarkSize] = useState(25);
  const [tabIndex, setTabIndex] = useState(0);
  const [isBulkUpload, setIsBulkUpload] = useState(false);
  const [bulkFiles, setBulkFiles] = useState<FileList | null>(null);
  const [bulkResults, setBulkResults] = useState<{ filename: string; data: string }[] | null>(null);
  const [isDownloadDisabled, setIsDownloadDisabled] = useState(false);

  useEffect(() => {
    const urlParams = new URLSearchParams(window.location.search);
    const status = urlParams.get('donation');
    if (status) {
      setDonationStatus(status);
      window.history.replaceState({}, document.title, window.location.pathname);
    }
  }, []);

  useEffect(() => {
    if (!user) {
      setIsDownloadDisabled(true);
    } else if (user.subscriptionStatus !== "active" && user.dailyDownloads >= 1) {
      setIsDownloadDisabled(true);
    } else {
      setIsDownloadDisabled(false);
    }
  }, [user]);

  const handleCloseSnackbar = () => {
    setDonationStatus(null);
  };

  const handleFileChange = useCallback((event: React.ChangeEvent<HTMLInputElement>) => {
    if (event.target.files && event.target.files[0]) {
      const selectedFile = event.target.files[0];
      setFile(selectedFile);
    }
  }, []);

  const handleBulkFileChange = useCallback((event: React.ChangeEvent<HTMLInputElement>) => {
    if (event.target.files && event.target.files.length > 0) {
      setBulkFiles(event.target.files);
    }
  }, []);

  const handleSubmit = useCallback(async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (!file && !bulkFiles) return;

    if (!user || !user.id) {
      setShowAuthModal(true);
      return;
    }

    setIsLoading(true);
    const formData = new FormData();
    formData.append('uniqueId', Date.now().toString());
    formData.append('opacity', opacity.toString());
    formData.append('spacing', spacing.toString());
    formData.append('userId', user.id);

    if (tabIndex === 0) {
      // Text watermark
      formData.append('text', watermarkText);
      formData.append('color', textColor);
      formData.append('fontSize', fontSize.toString());
    } else {
      // Image watermark
      if (watermarkImage) {
        formData.append('watermarkImage', watermarkImage);
      }
      formData.append('watermarkSize', watermarkSize.toString());
    }

    try {
      let response;
      if (isBulkUpload && bulkFiles) {
        // Bulk upload
        for (let i = 0; i < bulkFiles.length; i++) {
          formData.append('images', bulkFiles[i]);
        }
        response = await fetch('/api/watermark/bulk/' + (tabIndex === 0 ? 'text' : 'image'), {
          method: 'POST',
          body: formData,
        });
      } else if (file) {
        // Single file upload
        formData.append('image', file);
        response = await fetch('/api/watermark/' + (tabIndex === 0 ? 'text' : 'image'), {
          method: 'POST',
          body: formData,
        });
      } else {
        throw new Error('No file selected');
      }

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      if (data.results && Array.isArray(data.results)) {
        if (data.results.length === 1) {
          // Single image response
          setWatermarkedImage(data.results[0].data);
        } else {
          // Bulk image response
          setBulkResults(data.results);
        }
      } else {
        throw new Error('Unexpected response format');
      }
    } catch (error) {
      console.error('Error:', error);
      // Handle error (e.g., show error message to user)
    } finally {
      setIsLoading(false);
    }
  }, [file, bulkFiles, watermarkText, textColor, opacity, fontSize, spacing, tabIndex, watermarkImage, watermarkSize, user, isBulkUpload]);

  const handleDownload = useCallback(() => {
    if (isDownloadDisabled) {
      if (!user) {
        setShowAuthModal(true);
      } else {
        alert("Daily download limit reached. Please upgrade your subscription.");
      }
      return;
    }
    if (watermarkedImage) {
      // Create a temporary anchor element
      const link = document.createElement('a');
      link.href = watermarkedImage;
      link.download = 'watermarked_image.png';
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
    }
  }, [watermarkedImage, isDownloadDisabled, user]);

  const handleTabChange = (_event: React.SyntheticEvent, newValue: number) => {
    setTabIndex(newValue);
  };

  return (
    <ThemeProvider theme={theme}>
      <Elements stripe={stripePromise}>
        <Box sx={{ 
          display: 'flex', 
          flexDirection: 'column', 
          minHeight: '100vh', 
          bgcolor: 'transparent', 
          position: 'relative',
          backdropFilter: 'blur(10px)',
          backgroundColor: 'rgba(255, 255, 255, 0.1)',
        }}>
          <BackgroundBlobs />
          <AppBarComponent />
          <Container maxWidth="xl" sx={{ mt: 4, mb: 4, flexGrow: 1, pb: 4, display: 'flex', flexDirection: 'column' }}>
            <Routes>
              <Route path="/" element={
                <Grid container spacing={3}>
                  <Grid item xs={12} md={4}>
                    <Grid container direction="column" spacing={3}>
                      <Grid item>
                        <Paper elevation={3} sx={{ p: 3, backgroundColor: 'rgba(255, 255, 255, 0.1)', backdropFilter: 'blur(5px)', border: '1px solid rgba(255, 255, 255, 0.3)', borderRadius: '16px', boxShadow: '0 4px 30px rgba(0, 0, 0, 0.1)' }}>
                          <Typography variant="h6" gutterBottom color="secondary">Watermark Controls</Typography>
                          <Tabs
                            value={tabIndex}
                            onChange={handleTabChange}
                            aria-label="watermark control tabs"
                            variant="fullWidth"
                            sx={{
                              borderBottom: 1,
                              borderColor: 'divider',
                              '& .MuiTab-root': {
                                color: 'text.secondary',
                                '&.Mui-selected': {
                                  color: 'palette.primary.dark', // Green color
                                  borderBottom: '2px solid',
                                  borderColor: 'palette.primary.dark', // Green color
                                  fontWeight: 'bold'
                                },
                                '&:focus': {
                                  outline: 'none',
                                },
                              },
                            }}
                          >
                            <Tab label="Text Watermark" />
                            <Tab label="Image Watermark" />
                          </Tabs>
                          <Box sx={{ p: 3 }}>
                            <Box component="form" onSubmit={handleSubmit} sx={{ display: 'flex', flexDirection: 'column', gap: 1, bgcolor: 'background.paper', borderRadius: 2, p: 2 }}>
                              <Button variant="outlined" component="label" startIcon={<CloudUpload />} sx={{ mb: 1 }}>
                                {isBulkUpload ? 'Upload Bulk Images' : 'Upload Image'}
                                <input
                                  type="file"
                                  hidden
                                  onChange={isBulkUpload ? handleBulkFileChange : handleFileChange}
                                  accept="image/*"
                                  multiple={isBulkUpload}
                                />
                              </Button>
                              <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                                <FormControlLabel
                                  control={
                                    <Checkbox
                                      checked={isBulkUpload}
                                      onChange={(e) => setIsBulkUpload(e.target.checked)}
                                      disabled={user?.subscriptionStatus !== "active"}
                                    />
                                  }
                                  label="Bulk Upload"
                                />
                                <Chip
                                  label="PRO"
                                  size="small"
                                  sx={{
                                    ml: 1,
                                    background: 'linear-gradient(45deg, #3a86ff 30%, #023e8a 90%)',
                                    color: 'white',
                                    fontWeight: 'bold',
                                    fontSize: '0.7rem',
                                  }}
                                />
                              </Box>
                              {tabIndex === 0 && (
                                <>
                                  <TextField
                                    fullWidth
                                    label="Watermark Text"
                                    value={watermarkText}
                                    onChange={(e) => setWatermarkText(e.target.value)}
                                    sx={{ mb: 1 }}
                                  />
                                  <TextField
                                    fullWidth
                                    type="color"
                                    label="Text Color"
                                    value={textColor}
                                    onChange={(e) => setTextColor(e.target.value)}
                                    sx={{ mb: 1 }}
                                  />
                                  <Grid container spacing={2} sx={{ mb: 1 }}>
                                    <Grid item xs={4}>
                                      <Typography variant="body2">Opacity</Typography>
                                      <Slider
                                        value={opacity}
                                        onChange={(_, newValue) => setOpacity(newValue as number)}
                                        min={0}
                                        max={1}
                                        step={0.1}
                                        valueLabelDisplay="auto"
                                      />
                                    </Grid>
                                    <Grid item xs={4}>
                                      <Typography variant="body2">Font Size</Typography>
                                      <Slider
                                        value={fontSize}
                                        onChange={(_, newValue) => setFontSize(newValue as number)}
                                        min={10}
                                        max={200}
                                        step={1}
                                        valueLabelDisplay="auto"
                                      />
                                    </Grid>
                                    <Grid item xs={4}>
                                      <Typography variant="body2">Spacing</Typography>
                                      <Slider
                                        value={spacing}
                                        onChange={(_, newValue) => setSpacing(newValue as number)}
                                        min={50}
                                        max={1600}
                                        step={5}
                                        valueLabelDisplay="auto"
                                      />
                                    </Grid>
                                  </Grid>
                                </>
                              )}
                              {tabIndex === 1 && (
                                <>
                                  <Button variant="outlined" component="label" startIcon={<CloudUpload />} sx={{ mb: 1 }}>
                                    Choose Watermark Image
                                    <input type="file" hidden onChange={(e) => setWatermarkImage(e.target.files ? e.target.files[0] : null)} accept="image/*" />
                                  </Button>
                                  <Typography variant="body2" sx={{ mb: 1 }}>Watermark Size</Typography>
                                  <Slider
                                    value={watermarkSize}
                                    onChange={(_, newValue) => setWatermarkSize(newValue as number)}
                                    min={5}
                                    max={100}
                                    step={1}
                                    valueLabelDisplay="auto"
                                    sx={{ mb: 1 }}
                                  />
                                  <Grid container spacing={2} sx={{ mb: 1 }}>
                                    <Grid item xs={6}>
                                      <Typography variant="body2">Opacity</Typography>
                                      <Slider
                                        value={opacity}
                                        onChange={(_, newValue) => setOpacity(newValue as number)}
                                        min={0}
                                        max={1}
                                        step={0.1}
                                        valueLabelDisplay="auto"
                                      />
                                    </Grid>
                                    <Grid item xs={6}>
                                      <Typography variant="body2">Spacing</Typography>
                                      <Slider
                                        value={spacing}
                                        onChange={(_, newValue) => setSpacing(newValue as number)}
                                        min={50}
                                        max={1600}
                                        step={5}
                                        valueLabelDisplay="auto"
                                      />
                                    </Grid>
                                  </Grid>
                                </>
                              )}
                              <Button type="submit" variant="contained" color="primary" disabled={!file && !bulkFiles} sx={{ mt: 1 }}>
                                Apply Watermark
                              </Button>
                            </Box>
                          </Box>
                        </Paper>
                      </Grid>
                      <Grid item>
                        <Paper elevation={3} sx={{ p: 3, backgroundColor: 'rgba(255, 255, 255, 0.1)', backdropFilter: 'blur(5px)', border: '1px solid rgba(255, 255, 255, 0.3)', borderRadius: '16px', boxShadow: '0 4px 30px rgba(0, 0, 0, 0.1)' }}>
                          <Typography variant="h6" gutterBottom>Support the Project</Typography>
                          <DonationForm />
                        </Paper>
                      </Grid>
                    </Grid>
                  </Grid>
                  <Grid item xs={12} md={8}>
                    <Paper elevation={3} sx={{ p: 3, height: '100%', display: 'flex', flexDirection: 'column', backgroundColor: 'rgba(255, 255, 255, 0.1)', backdropFilter: 'blur(5px)', border: '1px solid rgba(255, 255, 255, 0.3)', borderRadius: '16px', boxShadow: '0 4px 30px rgba(0, 0, 0, 0.1)' }}>
                      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
                        <Typography variant="h6" gutterBottom color="secondary">Preview</Typography>
                        {watermarkedImage && (
                          <Button
                            variant="contained"
                            color="secondary"
                            startIcon={<DownloadIcon />}
                            onClick={handleDownload}
                            size="small"
                            disabled={isDownloadDisabled}
                          >
                            {isDownloadDisabled ? (user ? "Daily Limit Reached" : "Login to Download") : "Download Watermark"}
                          </Button>
                        )}
                      </Box>
                      <Box sx={{ flexGrow: 1, display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '400px' }}>
                        <ImagePreview
                          previews={
                            isBulkUpload 
                              ? bulkFiles ? Array.from(bulkFiles).map((file: File) => URL.createObjectURL(file)) : []
                              : file 
                                ? [URL.createObjectURL(file)] 
                                : []
                          }
                          isLoading={isLoading}
                          watermarkedImages={
                            isBulkUpload
                              ? (bulkResults || []).map(result => result.data)
                              : watermarkedImage
                                ? [watermarkedImage]
                                : []
                          }
                        />
                      </Box>
                    </Paper>
                  </Grid>
                </Grid>
              } />
              <Route path="/signin" element={<SignIn />} />
              <Route path="/signup" element={<SignUp />} />
              <Route path="/admin" element={<AdminPage />} />
              <Route path="/subscription" element={<Subscription />} />
              <Route path="/subscribe/success" element={<SubscriptionSuccess />} />
              <Route path="/subscribe/cancel" element={<SubscriptionCancel />} />
              <Route path="/settings" element={<SettingsPage />} />
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
                {'Copyright  '}
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
      <AuthModal
        open={showAuthModal}
        onClose={() => setShowAuthModal(false)}
      />
    </ThemeProvider>
  );
}

export default App;