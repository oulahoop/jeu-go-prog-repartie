# GUIDE 1

## 1 - Etude du code fourni

```
En étudiant le code, vous devez expliquer (par exemple au moyen de schémas) 
comment sedéroule un tour de la boucle principale du jeu 
(i.e. un appel à Update) : 

- Quelles sont les fonctions qui sont appelées ? 
- A quoi servent-elles ? 
- Comment sait-on à quelle étape du jeu on est ? 
- À quoi sert l’argument -tps lorsqu’on lance le jeu ?
```

Les fonctions du jeu sont qui sont appelées sont :
- **InitGame:** Cela initialise le jeu et les variables, il est appelé une fois avant le lancement.


- **Update:** Cela met à jour le jeu en fonction de son état (state), elle est appelé à chaque frame (60 fois par seconde)
    - Une fonction Update est appelée depuis la struct Game et retourne un booléen (bool), ce booléan sert à savoir si on change de state ou non. 
      Cette fonction appelera les fonctions de Update spécifique de l'état du jeu (state). 


- **Draw:** Cela dessine le jeu en fonction de son état (state), elle est appelé après Update afin de dessiner le jeu.
    -  Une fonction Draw est appelée depuis la struct Game et retrourne une image (ebiten). Cette fonction appelera les fonctions de Draw spécifique de l'état du jeu (state). 
Elles servent à créer le jeu et à le mettre à jour.

Dans les états on retrouve : 
- **WelcomeScreen :** C'est l'écran d'accueil du jeu 
- **ChooseRunner :** C'est l'écran où l'on choisit le personnage que l'on veut jouer
- **LaunchRun :** C'est l'écran où l'on lance la course (le 3,2,1) 
- **Run :** C'est l'écran où l'on joue
- **Result :** C'est la fin de la course, on voit le résultat de la course

On sait à quelle étape du jeu on est grâce à la variable state qui est incrémentée à chaque étape. Certains état on une variable d'état elle-même (

L'argument -tps lorsqu'on lance le jeu sert à afficher les FPS (Frame Per Second) et le temps de calcul de la frame. Il est limité à 60.
