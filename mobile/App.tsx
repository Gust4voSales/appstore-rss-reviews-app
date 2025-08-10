import { StatusBar } from 'expo-status-bar';

import './global.css';
import { SafeAreaProvider, SafeAreaView } from 'react-native-safe-area-context';
import { ReviewsScreen } from 'screens/ReviewsScreen';

export default function App() {
  return (
    <SafeAreaProvider>
      <SafeAreaView className="flex-1">
        <ReviewsScreen />
        <StatusBar style="auto" />
      </SafeAreaView>
    </SafeAreaProvider>
  );
}
