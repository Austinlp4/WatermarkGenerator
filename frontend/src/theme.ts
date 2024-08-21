import { createTheme, ThemeOptions } from '@mui/material/styles';

const themeOptions: ThemeOptions = {
  palette: {
    primary: {
      main: '#0A192F',
      light: '#172A45',
      dark: '#020C1B',
    },
    secondary: {
      main: '#64FFDA',
      light: '#A8FFF1',
      dark: '#00C9A7',
    },
    background: {
      default: '#0A192F',
      paper: '#1E2A45',
    },
    text: {
      primary: '#FFFFFF',
      secondary: '#B2BAC2',
      disabled: '#8892B0',
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
        root: () => ({
          borderRadius: 4,
          textTransform: 'none',
          fontSize: '0.875rem',
          padding: '10px 20px',
          boxShadow: '0 4px 6px rgba(0, 0, 0, 0.12)',
          transition: 'all 0.3s cubic-bezier(0.645, 0.045, 0.355, 1)',
          '&:hover': {
            transform: 'translateY(-2px)',
            boxShadow: '0 7px 14px rgba(0, 0, 0, 0.18)',
          },
          color: '#333333',
          '&.MuiButton-contained': {
            backgroundColor: '#64FFDA',
            color: '#333333',
            '&:hover': {
              backgroundColor: '#00C9A7',
            },
          },
          '&.MuiButton-outlined': {
            borderColor: '#64FFDA',
            color: '#64FFDA',
            '&:hover': {
              backgroundColor: 'rgba(100, 255, 218, 0.1)',
            },
          },
        }),
      },
    },
    MuiPaper: {
      styleOverrides: {
        root: {
          borderRadius: 8,
          boxShadow: '0 10px 30px rgba(0, 0, 0, 0.2)',
        },
      },
    },
    MuiAppBar: {
      styleOverrides: {
        root: {
          background: 'linear-gradient(90deg, #0A192F 0%, #172A45 100%)',
        },
      },
    },
    MuiOutlinedInput: {
      styleOverrides: {
        root: {
          '& .MuiOutlinedInput-notchedOutline': {
            borderColor: 'rgba(255, 255, 255, 0.3)',
          },
          '&:hover .MuiOutlinedInput-notchedOutline': {
            borderColor: 'rgba(255, 255, 255, 0.5)',
          },
          '&.Mui-focused .MuiOutlinedInput-notchedOutline': {
            borderColor: '#64FFDA',
          },
        },
      },
    },
    MuiInputLabel: {
      styleOverrides: {
        root: {
          color: '#B2BAC2',
          '&.Mui-focused': {
            color: '#64FFDA',
          },
        },
      },
    },
    MuiSlider: {
      styleOverrides: {
        root: {
          color: '#64FFDA',
        },
      },
    },
  },
};

const theme = createTheme(themeOptions);

export default theme;
