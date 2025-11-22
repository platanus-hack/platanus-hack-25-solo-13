import type { Subject } from '$lib/constants/subjects';

// Mission interface
export interface Mission {
  id: number;
  subject: string;
  title: string;
  time: string;
  reward: string;
  state: 'start' | 'progress' | 'done';
}

// Event interface
export interface Event {
  id: number;
  title: string;
  subtitle: string;
  color: string;
}

// Activity interface
export interface Activity {
  id: number;
  text: string;
  time: string;
  icon: string;
}

// Missions organized by category
export interface MissionsData {
  daily: Mission[];
  weekly: Mission[];
  story: Mission[];
  side: Mission[];
}

// Dashboard store state
class DashboardStore {
  missions = $state<MissionsData>({
    daily: [
      { id: 1, subject: 'Math', title: 'Linear Equations Practice', time: '15 min', reward: '50 XP', state: 'start' },
      { id: 2, subject: 'Language', title: 'Reading Comprehension', time: '10 min', reward: '30 XP', state: 'done' },
      { id: 3, subject: 'History', title: 'Chilean Independence', time: '20 min', reward: '60 XP', state: 'start' }
    ],
    weekly: [
      { id: 4, subject: 'Math', title: 'Complete Unit 3: Functions', time: '2 hrs', reward: '500 XP', state: 'progress' },
      { id: 5, subject: 'Physics', title: 'Lab Report: Motion', time: '45 min', reward: '200 XP', state: 'start' }
    ],
    story: [
      { id: 6, subject: 'Campaign', title: 'Chapter 2: The Algebra Realm', time: 'Ongoing', reward: 'Badge', state: 'progress' }
    ],
    side: [
      { id: 7, subject: 'Challenge', title: 'Speed Math: Mental Calculation', time: '5 min', reward: '10 SP', state: 'start' }
    ]
  });

  events = $state<Event[]>([
    { id: 1, title: 'Exam Sprint Week', subtitle: 'Boost your PAES score', color: 'from-indigo-600 to-violet-600' },
    { id: 2, title: 'Math Boss Challenge', subtitle: 'Beat the Function Dragon', color: 'from-rose-600 to-orange-600' }
  ]);

  activities = $state<Activity[]>([
    { id: 1, text: "You unlocked 'Equation Explorer'", time: '2h ago', icon: '🏆' },
    { id: 2, text: 'Teacher left feedback on essay', time: '4h ago', icon: '💬' },
    { id: 3, text: 'New PAES simulator available', time: '1d ago', icon: '🆕' }
  ]);

  subjects = $state<Subject[]>([]);

  // Derived: Count active missions (not done)
  get activeMissionCount(): number {
    return Object.values(this.missions)
      .flat()
      .filter(m => m.state !== 'done')
      .length;
  }

  // Derived: Count pending activities
  get activityCount(): number {
    return this.activities.length;
  }

  // Derived: Count active events
  get eventCount(): number {
    return this.events.length;
  }

  // Update missions
  updateMissions(missions: MissionsData) {
    this.missions = missions;
  }

  // Update events
  updateEvents(events: Event[]) {
    this.events = events;
  }

  // Update activities
  updateActivities(activities: Activity[]) {
    this.activities = activities;
  }

  // Update subjects
  updateSubjects(subjects: Subject[]) {
    this.subjects = subjects;
  }

  // Add activity
  addActivity(activity: Activity) {
    this.activities = [activity, ...this.activities];
  }

  // Mark mission as done
  markMissionDone(missionId: number) {
    // Find and update mission in any category
    for (const category of Object.keys(this.missions) as Array<keyof MissionsData>) {
      const missionIndex = this.missions[category].findIndex(m => m.id === missionId);
      if (missionIndex !== -1) {
        this.missions[category][missionIndex].state = 'done';
        break;
      }
    }
  }

  // Mark mission as in progress
  markMissionInProgress(missionId: number) {
    for (const category of Object.keys(this.missions) as Array<keyof MissionsData>) {
      const missionIndex = this.missions[category].findIndex(m => m.id === missionId);
      if (missionIndex !== -1) {
        this.missions[category][missionIndex].state = 'progress';
        break;
      }
    }
  }
}

// Create singleton instance
export const dashboardStore = new DashboardStore();
