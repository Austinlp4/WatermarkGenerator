import { createTheme, ThemeOptions } from '@mui/material/styles';

const themeOptions: ThemeOptions = {
  palette: {
    primary: {
      main: '#3a86ff',
      light: '#8ecae6',
      dark: '#023e8a',
    },
    secondary: {
      main: '#023e8a',
      light: '#ffafcc',
      dark: '#8f00ff',
    },
    background: {
      default: '#ffffff',
      paper: '#f8f9fa',
    },
    text: {
      primary: '#333333',
      secondary: '#6c757d',
      disabled: '#adb5bd',
    },
  },
  typography: {
    fontFamily: '"Inter", "Roboto", "Helvetica", "Arial", sans-serif',
    h4: {
      fontWeight: 700,
    },
  },
  components: {
    MuiButton: {
      styleOverrides: {
        root: {
          borderRadius: 8,
          textTransform: 'none',
          fontWeight: 600,
          boxShadow: '0 2px 4px rgba(0, 0, 0, 0.1)',
          '&:hover': {
            boxShadow: '0 4px 6px rgba(0, 0, 0, 0.15)',
          },
        },
        contained: {
          backgroundColor: '#3a86ff',
          color: '#ffffff',
          '&:hover': {
            backgroundColor: '#023e8a',
          },
        },
        outlined: {
          borderColor: '#3a86ff',
          color: '#3a86ff',
          '&:hover': {
            backgroundColor: 'rgba(58, 134, 255, 0.1)',
          },
        },
      },
    },
    MuiPaper: {
      styleOverrides: {
        root: {
          backgroundColor: '#ffffff',
          borderRadius: 16,
          boxShadow: '0 4px 6px rgba(0, 0, 0, 0.1)',
        },
      },
    },
    MuiAppBar: {
      styleOverrides: {
        root: {
          backgroundColor: '#ffffff',
          color: '#333333',
          boxShadow: '0 2px 4px rgba(0, 0, 0, 0.1)',
        },
      },
    },
    MuiOutlinedInput: {
      styleOverrides: {
        root: {
          '& .MuiOutlinedInput-notchedOutline': {
            borderColor: 'rgba(0, 0, 0, 0.23)',
          },
          '&:hover .MuiOutlinedInput-notchedOutline': {
            borderColor: 'rgba(0, 0, 0, 0.87)',
          },
          '&.Mui-focused .MuiOutlinedInput-notchedOutline': {
            borderColor: '#3a86ff',
          },
        },
      },
    },
    MuiInputLabel: {
      styleOverrides: {
        root: {
          color: '#6c757d',
          '&.Mui-focused': {
            color: '#3a86ff',
          },
        },
      },
    },
    MuiSlider: {
      styleOverrides: {
        root: {
          color: '#3a86ff',
        },
      },
    },
  },
};

const theme = createTheme(themeOptions);

export default theme;
