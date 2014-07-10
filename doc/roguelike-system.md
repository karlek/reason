Ritual system
=============

The player can choose a ritual to preform which will give a reward after it's 
completion. For example:

"Defeat 1000 agiri monsters with only a broom to achieve broom magic."

"Fast for 7 days to achieve conjure food."
	
	Stat loss ensues when character doesn't eat.
		Stat death: stats must always be higher than 0.
	Character needs to have status Engorged.
	Character needs ring of sustenence.


Stats
=====

Magic power
	
	Level / damage / hit chance with spells.

Endurance

	HP, MP, how long one can endure a ritual.

Speed

	Act more quickly.

Strength

	Improves physical damage.

Perception

	Increases area of field of view.
	Consider multiple senses.


Story
=====

You are a unranked black magician.

Kill Shion and Emma to retrieve the power of Echnida.
	
	Randomly start somewhere in the world.

Graphics
========

Standard
--------

	$ gold
	: magic book
	! potions
	% corpse
	= ring
	" amulet
	? scroll
	
	^ trap
	+-/ doors
	_ altar
	> stairs down
	< stairs up
	. soil

Monsters
--------

a-z (lowercase letters for common)
	
A-Z (capital letters for uncommon/legendary/epic/rare monsters)
	
	D dragon

Unused
------

	&*()`]}[{\/',;

Combat
======

http://roguebasin.roguelikedevelopment.org/index.php?title=Thoughts_on_Combat_Models

Ranged combat
-------------

http://roguebasin.roguelikedevelopment.org/index.php?title=Two-Key_Targeting

Races
=====

Message system
==============

Queue
-----

My turn
Print -> 
Not my turn
	event: add to queue
	event: add to queue
	event: add to queue
My turn
Print -> event, event, event!

Draw lines on grids
===================

http://roguebasin.roguelikedevelopment.org/index.php?title=Bresenham%27s_Line_Algorithm

Design decisions
================

http://roguebasin.roguelikedevelopment.org/index.php?title=Dungeon_persistence
http://roguebasin.roguelikedevelopment.org/index.php?title=RL_Terrain
http://roguebasin.roguelikedevelopment.org/index.php?title=Shop
http://roguebasin.roguelikedevelopment.org/index.php?title=Trap
http://roguebasin.roguelikedevelopment.org/index.php?title=Data_structures_for_the_map
http://roguebasin.roguelikedevelopment.org/index.php?title=Code_design_basics
http://roguebasin.roguelikedevelopment.org/index.php?title=The_Role_of_Hunger
http://roguebasin.roguelikedevelopment.org/index.php?title=Fractals
http://roguebasin.roguelikedevelopment.org/index.php?title=Items
http://roguebasin.roguelikedevelopment.org/index.php?title=Interesting_Critical_Hits
http://roguebasin.roguelikedevelopment.org/index.php?title=Info_Files
http://roguebasin.roguelikedevelopment.org/index.php?title=Info_File_Variant_-_Compile-to-Code
http://roguebasin.roguelikedevelopment.org/index.php?title=Monster_attacks
http://roguebasin.roguelikedevelopment.org/index.php?title=Monster_attacks_in_a_structured_list
http://roguebasin.roguelikedevelopment.org/index.php?title=Overworld
http://roguebasin.roguelikedevelopment.org/index.php?title=Power_Curve
http://roguebasin.roguelikedevelopment.org/index.php?title=Preferred_Key_Controls
http://roguebasin.roguelikedevelopment.org/index.php?title=Roguelike_Interface
http://roguebasin.roguelikedevelopment.org/index.php?title=Running_code
http://roguebasin.roguelikedevelopment.org/index.php?title=Save_Files
http://roguebasin.roguelikedevelopment.org/index.php?title=Scrolling_map
http://roguebasin.roguelikedevelopment.org/index.php?title=Roguelike_Dev_FAQ

Line of sight
=============

http://roguebasin.roguelikedevelopment.org/index.php?title=Restrictive_Precise_Angle_Shadowcasting
http://roguebasin.roguelikedevelopment.org/index.php?title=Line_of_sight
http://roguebasin.roguelikedevelopment.org/index.php?title=Field_of_Vision
http://roguebasin.roguelikedevelopment.org/index.php?title=Discussion:Field_of_Vision
http://roguebasin.roguelikedevelopment.org/index.php?title=FOV_using_recursive_shadowcasting_-_improved

Later
=====

AI
===

http://roguebasin.roguelikedevelopment.org/index.php?title=A_Better_Monster_AI
http://roguebasin.roguelikedevelopment.org/index.php?title=Anticipating_wall-following_pathfinder
http://roguebasin.roguelikedevelopment.org/index.php?title=Need_driven_AI
http://roguebasin.roguelikedevelopment.org/index.php?title=Simple_way_to_prevent_jams_of_monsters_with_A*
http://roguebasin.roguelikedevelopment.org/index.php?title=Roguelike_Intelligence

Dungeon Generation
==================

http://roguebasin.roguelikedevelopment.org/index.php?title=Basic_directional_dungeon_generation
http://roguebasin.roguelikedevelopment.org/index.php?title=Abstract_Dungeons
http://roguebasin.roguelikedevelopment.org/index.php?title=Cellular_Automata_Method_for_Generating_Random_Cave-Like_Levels
http://roguebasin.roguelikedevelopment.org/index.php?title=Dungeon-Building_Algorithm
http://roguebasin.roguelikedevelopment.org/index.php?title=Diffusion-limited_aggregation
http://roguebasin.roguelikedevelopment.org/index.php?title=Grid_Based_Dungeon_Generator
http://roguebasin.roguelikedevelopment.org/index.php?title=Maze_Generation
